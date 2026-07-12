package main

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	userConn, err := grpc.NewClient("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Не удалось подключиться к User Service: %v", err)
	}
	defer userConn.Close()

	router := gin.New()

	router.Use(gin.Logger(), gin.Recovery())
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "ok",
		})
	})

	api := router.Group("/api/v1")
	{
		_ = api
	}

	log.Println("API Gateway запущен на порту :3000")
	if err := router.Run(":3000"); err != nil {
		log.Fatalf("Ошибка при запуске сервера: %v", err)
	}
}
