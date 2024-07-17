package http

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

type Server struct {
	Router *mux.Router
	server *http.Server
}

func NewServer(address string, WRtimeout, idleTimeout int) *Server {
	r := mux.NewRouter()

	return &Server{
		Router: r,
		server: &http.Server{
			Addr:         address,
			WriteTimeout: time.Second * time.Duration(WRtimeout),
			ReadTimeout:  time.Second * time.Duration(WRtimeout),
			IdleTimeout:  time.Second * time.Duration(idleTimeout),
		},
	}
}

func (s *Server) Start() {
	log.Printf("Starting the server on %s", s.server.Addr)

	go func() {
		if err := s.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Printf("Failed to listen and serve: %s", err.Error())
		}
	}()
}

func (s *Server) Shutdown(ctx context.Context) error {
	log.Printf("Shutting down the server")

	if err := s.server.Shutdown(ctx); err != nil {
		return err
	}

	log.Println("Server shutdown complete")
	return nil
}
