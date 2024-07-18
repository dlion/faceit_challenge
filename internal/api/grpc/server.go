package grpc

import (
	"log"
	"net"

	"google.golang.org/grpc"
)

type Server struct {
	server *grpc.Server
}

func NewServer() *Server {
	return &Server{server: grpc.NewServer()}
}

func (s *Server) RegisterService(service *grpc.ServiceDesc, impl interface{}) {
	s.server.RegisterService(service, impl)
}

func (s *Server) Start(address string) {
	log.Print("Starting gRPC server on ", address)

	lis, err := net.Listen("tcp", address)
	if err != nil {
		log.Fatal("failed to listen", err)
	}

	go func() {
		if err := s.server.Serve(lis); err != nil {
			log.Fatal("failed to serve gRPC", err)
		}
	}()
}

func (s *Server) Shutdown() {
	log.Print("Gracefully shutting down gRPC server")
	s.server.GracefulStop()
}
