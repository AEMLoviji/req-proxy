package proxy_test

import (
	"errors"
	"io"
	"net/http"
	"os"
	"req-proxy/http_client"
	"strings"
	"testing"

	"req-proxy/proxy"

	"github.com/stretchr/testify/assert"
)

func TestMain(m *testing.M) {
	http_client.StartMockups()
	os.Exit(m.Run())
}

func TestForwardHttpClientError(t *testing.T) {
	http_client.FlushMockups()
	http_client.AddMockup(http_client.Mock{
		Url:        "https://google.com",
		HttpMethod: http.MethodGet,
		Response:   nil,
		Err:        errors.New("http client error"),
	})

	req := proxy.ProxyRequest{
		Method:  proxy.Get,
		Url:     "https://google.com",
		Headers: nil,
	}

	resp, err := req.Forward()

	assert.Nil(t, resp)
	assert.NotNil(t, err)
	assert.EqualValues(t, "http client error", err.Error())
}

func TestForwardSuccess(t *testing.T) {
	http_client.FlushMockups()
	http_client.AddMockup(http_client.Mock{
		Url:        "http://google.com",
		HttpMethod: http.MethodGet,
		Response: &http.Response{
			StatusCode:    http.StatusOK,
			ContentLength: 123,
			Header: http.Header{
				"x-header-any": []string{"sads"},
			},
			Body: io.NopCloser(strings.NewReader(`{"message": 1}`)),
		},
	})

	req := proxy.ProxyRequest{
		Method:  proxy.Get,
		Url:     "http://google.com",
		Headers: nil,
	}

	resp, err := req.Forward()

	assert.NotNil(t, resp)
	assert.Nil(t, err)
	assert.EqualValues(t, http.StatusOK, resp.Status)
	assert.EqualValues(t, 123, resp.Length)
	assert.EqualValues(t, http.Header{"x-header-any": []string{"sads"}}, resp.Headers)
}

func TestForwardSuccessRequestIsTracked(t *testing.T) {
	http_client.FlushMockups()
	http_client.AddMockup(http_client.Mock{
		Url:        "http://google.com",
		HttpMethod: http.MethodGet,
		Response: &http.Response{
			StatusCode:    http.StatusOK,
			ContentLength: 123,
			Header:        nil,
			Body:          io.NopCloser(strings.NewReader(`{"message": 1}`)),
		},
	})

	req := proxy.ProxyRequest{
		Method:  proxy.Get,
		Url:     "http://google.com",
		Headers: nil,
	}

	resp, err := req.Forward()

	assert.NotNil(t, resp)
	assert.Nil(t, err)

	tl := proxy.TrackList()
	assert.EqualValues(t, 1, len(tl))
}
