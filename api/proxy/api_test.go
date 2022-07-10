package proxy_test

import (
	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"req-proxy/api/proxy"
	"req-proxy/domain"
	"req-proxy/observer"
	"testing"

	"github.com/google/uuid"
)

var (
	proxyService   domain.MockProxyService
	requestTracker observer.MockRequestHistoryTracker
	ts             *httptest.Server
)

func TestMain(m *testing.M) {
	mux := http.NewServeMux()

	pr := proxy.NewProxyResource(&proxyService, &requestTracker)
	pr.MapRoutes(mux)

	ts = httptest.NewServer(mux)

	code := m.Run()
	ts.Close()
	os.Exit(code)
}

func TestProxyRequest(t *testing.T) {
	proxyService.ForwardFunc = func(pr *domain.ProxyRequest) (*domain.ProxyResponse, error) {
		return &domain.ProxyResponse{}, nil
	}

	tests := []struct {
		name   string
		method domain.HttpMethod
		url    string
		status int
	}{
		{"google", domain.Get, "https://google.com", http.StatusOK},
		{"yandex", domain.Get, "https://yandex.com", http.StatusOK},
	}

	t.Parallel()

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			req, err := encode(&domain.ProxyRequest{
				Method:  tc.method,
				Url:     tc.url,
				Headers: nil,
			})
			if err != nil {
				t.Fatal("failed to encode request body")
			}

			res, _ := testRequest(t, ts, "GET", "/proxy", req)

			if res.StatusCode != tc.status {
				t.Errorf("got http status %d, want: %d", res.StatusCode, tc.status)
			}
		})
	}
}

func TestProxyHistory(t *testing.T) {
	requestTrackerItem := map[uuid.UUID]observer.Entry{
		uuid.New(): {
			ClientRequest:      "fake client request",
			ThirdPartyResponse: "fake 3-rd party response",
		},
	}

	responseJson, _ := encode(requestTrackerItem)

	requestTracker.ListEntriesFunc = func() map[uuid.UUID]observer.Entry {
		return requestTrackerItem
	}

	res, resBody := testRequest(t, ts, "GET", "/proxy/history", nil)

	if res.StatusCode != http.StatusOK {
		t.Errorf("got http status %d, want: %d", res.StatusCode, http.StatusOK)
	}

	if resBody != responseJson.String() {
		t.Errorf("got incoreect response %s, want: %s", resBody, responseJson)
	}
}

func testRequest(t *testing.T, ts *httptest.Server, method, path string, body io.Reader) (*http.Response, string) {
	req, err := http.NewRequest(method, ts.URL+path, body)
	if err != nil {
		t.Fatal(err)
		return nil, ""
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
		return nil, ""
	}

	defer resp.Body.Close()

	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
		return nil, ""
	}

	return resp, string(respBody)
}

func encode(v interface{}) (*bytes.Buffer, error) {
	data := new(bytes.Buffer)
	err := json.NewEncoder(data).Encode(v)
	return data, err
}
