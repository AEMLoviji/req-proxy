package domain

import (
	"net/http"
	"req-proxy/http_client"
	"req-proxy/observer"

	"github.com/google/uuid"
)

type HttpMethod string

const (
	Get    HttpMethod = http.MethodGet
	Post              = http.MethodPost
	Put               = http.MethodPut
	Delete            = http.MethodDelete
)

type ProxyRequest struct {
	Method  HttpMethod          `json:"method"`
	Url     string              `json:"url"`
	Headers map[string][]string `json:"headers"`
}

type ProxyResponse struct {
	Id      uuid.UUID           `json:"id"`
	Status  int                 `json:"status"`
	Headers map[string][]string `json:"headers"`
	Length  int64               `json:"length"`
}

type ProxyService struct {
	RequestTracker observer.RequestTracker
}

func NewProxyService(rt observer.RequestTracker) *ProxyService {
	return &ProxyService{RequestTracker: rt}
}

func (ps *ProxyService) Forward(pr *ProxyRequest) (*ProxyResponse, error) {
	var proxyResponse *ProxyResponse

	defer func() {
		ps.RequestTracker.AddEntry(observer.Entry{
			ClientRequest:      pr,
			ThirdPartyResponse: proxyResponse,
		})
	}()

	res, err := http_client.Invoke(string(pr.Method), pr.Url, pr.Headers)
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()

	proxyResponse = &ProxyResponse{
		Id:      uuid.New(),
		Status:  res.StatusCode,
		Headers: res.Header,
		Length:  res.ContentLength,
	}

	return proxyResponse, nil
}
