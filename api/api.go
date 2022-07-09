package api

import (
	"net/http"

	proxyAPI "req-proxy/api/proxy"
	"req-proxy/domain"
	"req-proxy/observer"
)

func NewApi() *http.ServeMux {
	mux := http.NewServeMux()

	// api per resource. In scope of task we have only one resource called proxy
	rt := observer.NewProxyRequestTracker()
	ps := domain.NewProxyService(rt)
	pr := proxyAPI.NewProxyResource(ps, rt)
	pr.MapRoutes(mux)

	// ping<->pong logic to check if api can accept request
	mux.HandleFunc("/ping", func(rw http.ResponseWriter, r *http.Request) {
		rw.Header().Set("Content-Type", "text/plain")
		rw.WriteHeader(http.StatusOK)
		rw.Write([]byte("pong"))
	})

	return mux
}
