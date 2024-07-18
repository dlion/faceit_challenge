package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/dlion/faceit_challenge/internal/api/grpc"
	"github.com/dlion/faceit_challenge/internal/api/http"
	"github.com/dlion/faceit_challenge/internal/domain/services/user"
	repositories "github.com/dlion/faceit_challenge/internal/repositories/mongo"
	"github.com/dlion/faceit_challenge/pkg/proto/proto"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	WR_TIMEOUT   = 15
	IDLE_TIMEOUT = 60
)

func main() {
	httpServer := http.NewServer(":80", WR_TIMEOUT, IDLE_TIMEOUT)

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	clientOptions := options.Client().ApplyURI("mongodb://mongo:mongo@mongodb:27017/faceit?authSource=admin")
	mongoClient, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}

	err = mongoClient.Ping(ctx, nil)
	if err != nil {
		log.Fatalf("Failed to ping MongoDB: %v", err)
	}

	userRepo := repositories.NewUserRepositoryMongoImpl(mongoClient)
	userService := user.NewUserService(userRepo)

	grpcServer := grpc.NewServer()
	grpcUserHandler := grpc.NewUserGrpcHandler(userService)
	proto.RegisterUserServiceServer(grpcServer, grpcUserHandler)

	grpcServer.Start(":8080")

	userHandler := &http.UserHandler{UserService: userService}

	httpServer.Router.HandleFunc("/api/health", http.HealthCheckHandler).Methods("GET")
	httpServer.Router.HandleFunc("/api/users", userHandler.GetUsersHandler).Methods("GET")
	httpServer.Router.HandleFunc("/api/user", userHandler.AddUserHandler).Methods("POST")
	httpServer.Router.HandleFunc("/api/user/{id}", userHandler.UpdateUserHandler).Methods("PUT")
	httpServer.Router.HandleFunc("/api/user/{id}", userHandler.RemoveUserHandler).Methods("DELETE")

	httpServer.HttpServer.Handler = httpServer.Router

	httpServer.Start()

	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGTERM, syscall.SIGINT)

	sig := <-c
	log.Printf("Received signal: %s. Initiating graceful shutdown...", sig)

	grpcServer.Shutdown()
	err = httpServer.Shutdown(ctx)
	if err != nil {
		log.Fatalf("Failed to shutdown the HTTP server: %s", err.Error())
	}

	log.Println("Server gracefully stopped")
}
