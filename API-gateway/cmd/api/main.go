package main

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	userv1 "github.com/roman4k-gg/myGarden/pkg/user_v1"
)

func main() {
	userConn, err := grpc.NewClient("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Не удалось подключиться к User Service: %v", err)
	}
	defer userConn.Close()

	userClient := userv1.NewUserServiceClient(userConn)

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
	}

	log.Println("API Gateway запущен на порту :3000")
	if err := router.Run(":3000"); err != nil {
		log.Fatalf("Ошибка при запуске сервера: %v", err)
	}
}
