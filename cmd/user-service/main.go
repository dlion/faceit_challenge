package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/dlion/faceit_challenge/internal/api/http"
	"github.com/dlion/faceit_challenge/internal/domain/services/user"
	repositories "github.com/dlion/faceit_challenge/internal/repositories/mongo"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	WR_TIMEOUT   = 15
	IDLE_TIMEOUT = 60
)

func main() {
	httpServer := http.NewServer(":80", WR_TIMEOUT, IDLE_TIMEOUT)
	httpServer.Router.HandleFunc("/api/health", http.HealthCheckHandler).Methods("GET")

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	clientOptions := options.Client().ApplyURI("mongodb://mongo:mongo@localhost:27017/faceit?authSource=admin")
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
	userHandler := &http.UserHandler{UserService: userService}

	httpServer.Router.HandleFunc("/api/health", http.HealthCheckHandler).Methods("GET")
	httpServer.Router.HandleFunc("/api/user", userHandler.AddUserHandler).Methods("POST")

	httpServer.HttpServer.Handler = httpServer.Router

	httpServer.Start()

	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGTERM, syscall.SIGINT)

	sig := <-c
	log.Printf("Received signal: %s. Initiating graceful shutdown...", sig)

	err = httpServer.Shutdown(ctx)
	if err != nil {
		log.Fatalf("Failed to shutdown the HTTP server: %s", err.Error())
	}

	log.Println("Server gracefully stopped")
}
