package domain_test

import (
	"errors"
	"io"
	"net/http"
	"os"
	"req-proxy/domain"
	"req-proxy/http_client"
	"req-proxy/observer"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMain(m *testing.M) {
	http_client.StartMockups()
	c := m.Run()
	os.Exit(c)
}

func TestForwardHttpClientError(t *testing.T) {
	http_client.FlushMockups()
	http_client.AddMockup(http_client.Mock{
		Url:        "https://google.com",
		HttpMethod: http.MethodGet,
		Response:   nil,
		Err:        errors.New("http client error"),
	})

	req := domain.ProxyRequest{
		Method:  domain.Get,
		Url:     "https://google.com",
		Headers: nil,
	}

	rt := observer.NewProxyRequestTracker()
	ps := domain.NewProxyService(rt)
	resp, err := ps.Forward(&req)

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

	req := domain.ProxyRequest{
		Method:  domain.Get,
		Url:     "http://google.com",
		Headers: nil,
	}

	rt := observer.NewProxyRequestTracker()
	ps := domain.NewProxyService(rt)
	resp, err := ps.Forward(&req)

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

	req := domain.ProxyRequest{
		Method:  domain.Get,
		Url:     "http://google.com",
		Headers: nil,
	}

	rt := observer.NewProxyRequestTracker()
	ps := domain.NewProxyService(rt)
	resp, err := ps.Forward(&req)

	assert.NotNil(t, resp)
	assert.Nil(t, err)

	te := ps.RequestTracker.ListEntries()
	assert.EqualValues(t, 1, len(te))
}
