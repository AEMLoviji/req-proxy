package app

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"req-proxy/proxy"
)

func mapRoutes() {
	http.HandleFunc("/proxy", proxyHandler)
	http.HandleFunc("/proxy/logs", proxyLogsHandler)
}

func proxyHandler(rw http.ResponseWriter, r *http.Request) {
	var pr proxy.ProxyRequest

	if err := json.NewDecoder(r.Body).Decode(&pr); err != nil {
		log.Print(err.Error())
		replyError(rw, http.StatusBadRequest, "invalid request was sent")
		return
	}

	res, err := pr.Forward()
	if err != nil {
		log.Print(err.Error())
		replyError(rw, http.StatusInternalServerError, "error occured while proxying the request")
		return
	}

	replyJson(rw, res)
}

func proxyLogsHandler(rw http.ResponseWriter, r *http.Request) {
	replyJson(rw, proxy.TrackList())
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
