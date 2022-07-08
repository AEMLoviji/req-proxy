package app

import (
	"fmt"
	"log"
	"net/http"
)

func Start(port int) {
	mapRoutes()

	log.Printf("App runs on port %d", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), nil))
}
