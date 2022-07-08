package main

import (
	"flag"
	"req-proxy/app"
)

func main() {
	port := flag.Int("port", 8080, "Application port")
	flag.Parse()

	app.Start(*port)
}
