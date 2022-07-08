package http_client

import (
	"errors"
	"fmt"
	"net/http"
)

var (
	enabledMocks = false
	mocks        = make(map[string]*Mock)
)

type Mock struct {
	Url        string
	HttpMethod string
	Response   *http.Response
	Err        error
}

func buildMockId(httpMethod string, url string) string {
	return fmt.Sprintf("%s_%s", httpMethod, url)
}

func StartMockups() {
	enabledMocks = true
}

func FlushMockups() {
	mocks = make(map[string]*Mock)
}

func StopMockups() {
	enabledMocks = false
}

func AddMockup(mock Mock) {
	mocks[buildMockId(mock.HttpMethod, mock.Url)] = &mock
}

func Invoke(method string, url string, headers http.Header) (*http.Response, error) {
	if enabledMocks {
		return runMockHandler(method, url)
	}

	return invokeInternal(method, url, headers)
}

func runMockHandler(method, url string) (*http.Response, error) {
	mock := mocks[buildMockId(method, url)]
	if mock == nil {
		return nil, errors.New("no mockup found for given request")
	}

	return mock.Response, mock.Err
}

func invokeInternal(method string, url string, headers http.Header) (*http.Response, error) {
	request, _ := http.NewRequest(http.MethodPost, url, nil)
	request.Header = headers

	client := http.Client{}

	return client.Do(request)
}
