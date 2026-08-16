package main

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/roman4k-gg/myGarden/API-gateway/internal/config"
	"github.com/roman4k-gg/myGarden/API-gateway/internal/delivery"
	catalogv1 "github.com/roman4k-gg/myGarden/pkg/catalog_v1"
	userv1 "github.com/roman4k-gg/myGarden/pkg/user_v1"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("config error: %v", err)
	}

	userConn, err := grpc.NewClient(cfg.UserServiceAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("failed to connect to User Service: %v", err)
	}
	defer userConn.Close()

	catalogConn, err := grpc.NewClient(cfg.CatalogServiceAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("failed to connect to Catalog Service: %v", err)
	}
	defer catalogConn.Close()

	userClient := userv1.NewUserServiceClient(userConn)
	catalogClient := catalogv1.NewCatalogServiceClient(catalogConn)

	router := gin.New()
	router.Use(gin.Logger(), gin.Recovery())

	h := delivery.NewHandler(userClient, catalogClient, cfg.JWTSecret)
	h.InitRoutes(router)

	srv := &http.Server{
		Addr:    cfg.HTTPPort,
		Handler: router,
	}

	log.Printf("API Gateway listening on %s", cfg.HTTPPort)
	if err := srv.ListenAndServe(); err != nil {
		log.Fatalf("server error: %v", err)
	}
}
