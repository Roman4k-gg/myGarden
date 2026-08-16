package delivery

import (
	"net/http"

	"github.com/gin-gonic/gin"
	catalogv1 "github.com/roman4k-gg/myGarden/pkg/catalog_v1"
	userv1 "github.com/roman4k-gg/myGarden/pkg/user_v1"
)

type Handler struct {
	userClient    userv1.UserServiceClient
	catalogClient catalogv1.CatalogServiceClient
	jwtSecret     string
}

func NewHandler(uc userv1.UserServiceClient, cc catalogv1.CatalogServiceClient, jwtSecret string) *Handler {
	return &Handler{
		userClient:    uc,
		catalogClient: cc,
		jwtSecret:     jwtSecret,
	}
}

func (h *Handler) InitRoutes(router *gin.Engine) {
	api := router.Group("/api/v1")
	{
		api.POST("/register", h.register)
		api.POST("/login", h.login)

		protected := api.Group("/", AuthMiddleware(h.jwtSecret))
		{
			protected.GET("/plants", h.getPlants)
			protected.POST("/favorites", h.addFavorite)
			protected.GET("/favorites", h.getFavorites)
		}
	}
}

func (h *Handler) register(c *gin.Context) {
	var reqData struct {
		Email    string `json:"email" binding:"required"`
		Password string `json:"password" binding:"required"`
		Name     string `json:"name"`
	}
	if err := c.ShouldBindJSON(&reqData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}
	resp, err := h.userClient.Register(c.Request.Context(), &userv1.RegisterRequest{
		Email:    reqData.Email,
		Password: reqData.Password,
		Name:     reqData.Name,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, resp)
}

func (h *Handler) login(c *gin.Context) {
	var reqData struct {
		Email    string `json:"email" binding:"required"`
		Password string `json:"password" binding:"required"`
	}
	if err := c.ShouldBindJSON(&reqData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}
	resp, err := h.userClient.Login(c.Request.Context(), &userv1.LoginRequest{
		Email:    reqData.Email,
		Password: reqData.Password,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, resp)
}

func (h *Handler) getPlants(c *gin.Context) {
	resp, err := h.catalogClient.ListPlants(c.Request.Context(), &catalogv1.ListPlantsRequest{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, resp)
}

func (h *Handler) addFavorite(c *gin.Context) {
	userIDRaw, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	userIDFloat, ok := userIDRaw.(float64)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid user_id claim"})
		return
	}

	var reqData struct {
		PlantID int32 `json:"plant_id" binding:"required"`
	}
	if err := c.ShouldBindJSON(&reqData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	resp, err := h.catalogClient.AddFavorite(c.Request.Context(), &catalogv1.AddFavoriteRequest{
		UserId:  int32(userIDFloat),
		PlantId: reqData.PlantID,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, resp)
}

func (h *Handler) getFavorites(c *gin.Context) {
	userIDRaw, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	userIDFloat, ok := userIDRaw.(float64)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid user_id claim"})
		return
	}

	resp, err := h.catalogClient.GetFavorites(c.Request.Context(), &catalogv1.GetFavoritesRequest{
		UserId: int32(userIDFloat),
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, resp)
}