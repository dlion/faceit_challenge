package http

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

type Server struct {
	Router     *mux.Router
	HttpServer *http.Server
}

func NewServer(address string, WRtimeout, idleTimeout int) *Server {
	r := mux.NewRouter()

	return &Server{
		Router: r,
		HttpServer: &http.Server{
			Addr:         address,
			WriteTimeout: time.Second * time.Duration(WRtimeout),
			ReadTimeout:  time.Second * time.Duration(WRtimeout),
			IdleTimeout:  time.Second * time.Duration(idleTimeout),
		},
	}
}

func (s *Server) Start() {
	log.Printf("Starting the server on %s", s.HttpServer.Addr)

	go func() {
		if err := s.HttpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Printf("Failed to listen and serve: %s", err.Error())
		}
	}()
}

func (s *Server) Shutdown(ctx context.Context) error {
	log.Printf("Shutting down the server")

	if err := s.HttpServer.Shutdown(ctx); err != nil {
		return err
	}

	log.Println("Server shutdown complete")
	return nil
}
