package main

import (
	"log"
	"net/http"

	"strconv"

	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	catalogv1 "github.com/roman4k-gg/myGarden/pkg/catalog_v1"
	userv1 "github.com/roman4k-gg/myGarden/pkg/user_v1"
)

func main() {
	userConn, err := grpc.NewClient("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Не удалось подключиться к User Service: %v", err)
	}
	defer userConn.Close()

	userClient := userv1.NewUserServiceClient(userConn)

	catalogConn, err := grpc.NewClient("localhost:50052", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("failed to connect to Catalog Service: %v", err)
	}
	defer catalogConn.Close()
	catalogClient := catalogv1.NewCatalogServiceClient(catalogConn)

	router := gin.New()
	router.Use(gin.Logger(), gin.Recovery())

	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	api := router.Group("/api/v1")
	{
		api.POST("/register", func(c *gin.Context) {
			req := &userv1.RegisterRequest{
				Email:    "test@example.com",
				Password: "password123",
				Name:     "Test User",
			}

			resp, err := userClient.Register(c.Request.Context(), req)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}

			c.JSON(http.StatusOK, resp)
		})

		api.POST("/login", func(c *gin.Context) {
			req := &userv1.LoginRequest{
				Email: "test@example.com",
				Password: "password123",
			}

			resp, err := userClient.Login(c.Request.Context(), req)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
			c.JSON(http.StatusOK, resp)
		})

		api.GET("/plants", func(c *gin.Context) {
			resp, err := catalogClient.ListPlants(c.Request.Context(), &catalogv1.ListPlantsRequest{})
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
			c.JSON(http.StatusOK, resp)
		})

		api.POST("/favorites", func(c *gin.Context) {
			userID, _ := strconv.Atoi(c.DefaultQuery("user_id", "1"))
			plantID, _ := strconv.Atoi(c.DefaultQuery("plant_id", "1"))
			
			req := &catalogv1.AddFavoriteRequest{
				UserId:  int32(userID),
				PlantId: int32(plantID),
			}
			resp, err := catalogClient.AddFavorite(c.Request.Context(), req)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
			c.JSON(http.StatusOK, resp)
		})

		api.GET("/favorites", func(c *gin.Context) {
			userID, _ := strconv.Atoi(c.DefaultQuery("user_id", "1"))
			req := &catalogv1.GetFavoritesRequest{UserId: int32(userID)}
			
			resp, err := catalogClient.GetFavorites(c.Request.Context(), req)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
			c.JSON(http.StatusOK, resp)
		})
	}

	log.Println("API Gateway запущен на порту :3000")
	if err := router.Run(":3000"); err != nil {
		log.Fatalf("Ошибка при запуске сервера: %v", err)
	}
}
