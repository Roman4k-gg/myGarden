package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	catalogv1 "github.com/roman4k-gg/myGarden/pkg/catalog_v1"
	userv1 "github.com/roman4k-gg/myGarden/pkg/user_v1"
)

func authMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header is required"})
			c.Abort()
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header format must be Bearer {token}"})
			c.Abort()
			return
		}

		tokenString := parts[1]
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return []byte("my-super-secret-key"), nil
		})

		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			c.Set("user_id", claims["user_id"])
		} else {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token claims"})
			c.Abort()
			return
		}

		c.Next()
	}
}

func main() {
	userServiceAddr := os.Getenv("USER_SERVICE_ADDR")
	if userServiceAddr == "" {
		userServiceAddr = "localhost:50051"
	}
	userConn, err := grpc.NewClient(userServiceAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("failed to connect to User Service: %v", err)
	}
	defer userConn.Close()

	userClient := userv1.NewUserServiceClient(userConn)

	catalogServiceAddr := os.Getenv("CATALOG_SERVICE_ADDR")
	if catalogServiceAddr == "" {
		catalogServiceAddr = "localhost:50052"
	}
	catalogConn, err := grpc.NewClient(catalogServiceAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
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

		protected := api.Group("/")
		protected.Use(authMiddleware())
		{
			protected.GET("/plants", func(c *gin.Context) {
				resp, err := catalogClient.ListPlants(c.Request.Context(), &catalogv1.ListPlantsRequest{})
				if err != nil {
					c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
					return
				}
				c.JSON(http.StatusOK, resp)
			})

			protected.POST("/favorites", func(c *gin.Context) {
				userIDStr := fmt.Sprintf("%v", c.MustGet("user_id"))
				userID, _ := strconv.Atoi(userIDStr)
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

			protected.GET("/favorites", func(c *gin.Context) {
				userIDStr := fmt.Sprintf("%v", c.MustGet("user_id"))
				userID, _ := strconv.Atoi(userIDStr)
				req := &catalogv1.GetFavoritesRequest{UserId: int32(userID)}
				
				resp, err := catalogClient.GetFavorites(c.Request.Context(), req)
				if err != nil {
					c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
					return
				}
				c.JSON(http.StatusOK, resp)
			})
		}
	}

	log.Println("API Gateway запущен на порту :3000")
	if err := router.Run(":3000"); err != nil {
		log.Fatalf("Ошибка при запуске сервера: %v", err)
	}
}
