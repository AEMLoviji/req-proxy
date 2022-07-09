package api

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
)

type Server struct {
	*http.Server
}

func NewServer(port int) *Server {
	log.Println("configuring server...")
	api := NewApi()

	srv := http.Server{
		Addr:    fmt.Sprintf(":%d", port),
		Handler: api,
	}

	return &Server{&srv}
}

func (srv *Server) Start() {
	log.Println("starting server...")
	go func() {
		if err := srv.ListenAndServe(); err != http.ErrServerClosed {
			panic(err)
		}
	}()
	log.Printf("Listening on %s\n", srv.Addr)

	// graceful shutdown logic.
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	sig := <-quit
	log.Println("Shutting down server... Reason:", sig)

	if err := srv.Shutdown(context.Background()); err != nil {
		panic(err)
	}
	log.Println("Server gracefully stopped")
}
