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
	"github.com/dlion/faceit_challenge/internal/api/http/handlers"
	"github.com/dlion/faceit_challenge/internal/domain/services/user"
	repositories "github.com/dlion/faceit_challenge/internal/repositories/mongo"
	"github.com/dlion/faceit_challenge/pkg/notifier"
	"github.com/dlion/faceit_challenge/pkg/proto"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	WR_TIMEOUT      = 15
	IDLE_TIMEOUT    = 60
	MONGODB_ENV_VAR = "MONGODB_URI"
)

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	mongodbURI := os.Getenv(MONGODB_ENV_VAR)
	if mongodbURI == "" {
		log.Fatalf("%s environment variable is not set", MONGODB_ENV_VAR)
	}

	mongoClient := createMongoClient(ctx, mongodbURI)
	userRepo := repositories.NewUserRepositoryMongoImpl(mongoClient)
	userChangeNotifier := notifier.NewNotifier()
	userService := user.NewUserService(userRepo, &userChangeNotifier)

	grpcServer := createGrpcServer(userService)
	grpcServer.Start(":8080")

	healthcheckHandler := handlers.NewHealthCheckHandler(mongoClient)
	userHandler := handlers.NewUserHandler(userService)

	httpServer := defineHandlers(healthcheckHandler, userHandler)
	httpServer.Start()

	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGTERM, syscall.SIGINT)

	sig := <-c
	log.Printf("Received signal: %s. Initiating graceful shutdown...", sig)

	shutdownServers(ctx, grpcServer, httpServer)

	log.Println("Server gracefully stopped")
}

func shutdownServers(ctx context.Context, grpcServer *grpc.Server, httpServer *http.Server) {
	grpcServer.Shutdown()
	err := httpServer.Shutdown(ctx)
	if err != nil {
		log.Fatalf("Failed to shutdown the HTTP server: %s", err.Error())
	}
}

func defineHandlers(healthcheck *handlers.HealthCheckHandler, user *handlers.UserHandler) *http.Server {
	httpServer := http.NewServer(":80", WR_TIMEOUT, IDLE_TIMEOUT)

	httpServer.Router.HandleFunc("/api/health", healthcheck.HealthCheckHandler).Methods("GET")
	httpServer.Router.HandleFunc("/api/users", user.GetUsersHandler).Methods("GET")
	httpServer.Router.HandleFunc("/api/user", user.AddUserHandler).Methods("POST")
	httpServer.Router.HandleFunc("/api/user/{id}", user.UpdateUserHandler).Methods("PUT")
	httpServer.Router.HandleFunc("/api/user/{id}", user.RemoveUserHandler).Methods("DELETE")
	httpServer.HttpServer.Handler = httpServer.Router

	return httpServer
}

func createGrpcServer(userService *user.UserServiceImpl) *grpc.Server {
	grpcServer := grpc.NewServer()
	grpcUserHandler := grpc.NewUserGrpcHandler(userService)
	proto.RegisterUserServiceServer(grpcServer, grpcUserHandler)
	return grpcServer
}

func createMongoClient(ctx context.Context, mongodbURI string) *mongo.Client {
	clientOptions := options.Client().ApplyURI(mongodbURI)
	mongoClient, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}
	return mongoClient
}
