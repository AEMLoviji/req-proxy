package proxy

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"req-proxy/domain"
	"req-proxy/observer"
)

type ProxyResource struct {
	ProxySvc       domain.ProxyServiceInterface
	RequestTracker observer.RequestHistoryTracker
}

func NewProxyResource(ps domain.ProxyServiceInterface, rt observer.RequestHistoryTracker) *ProxyResource {
	return &ProxyResource{
		ProxySvc:       ps,
		RequestTracker: rt,
	}
}

func (p *ProxyResource) MapRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/proxy", p.proxyHandler)
	mux.HandleFunc("/proxy/history", p.proxyHistoryHandler)
}

func (p *ProxyResource) proxyHandler(rw http.ResponseWriter, r *http.Request) {
	var req domain.ProxyRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Print(err.Error())
		replyError(rw, http.StatusBadRequest, "invalid request was sent")
		return
	}

	res, err := p.ProxySvc.Forward(&req)
	if err != nil {
		log.Print(err.Error())
		replyError(rw, http.StatusInternalServerError, "error occured while proxying the request")
		return
	}

	replyJson(rw, res)
}

func (p *ProxyResource) proxyHistoryHandler(rw http.ResponseWriter, r *http.Request) {
	replyJson(rw, p.RequestTracker.ListEntries())
}

func replyError(rw http.ResponseWriter, status int, format string, args ...interface{}) {
	http.Error(rw, fmt.Sprintf(format, args...), status)
}

func replyJson(rw http.ResponseWriter, model interface{}) {
	rw.Header().Set("Content-Type", "application/json; charset=utf-8")
	rw.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(rw).Encode(model); err != nil {
		log.Print(err.Error())
	}
}
