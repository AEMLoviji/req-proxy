package main

import (
	"flag"
	"req-proxy/api"
)

func main() {
	port := flag.Int("port", 8080, "Application port")
	flag.Parse()

	server := api.NewServer(*port)
	server.Start()
}
