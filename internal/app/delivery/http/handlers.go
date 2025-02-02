package http

import (
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	"wallet-app/docs"

	"wallet-app/internal/app/services"
)

type Handler struct {
	services *services.Service
}

func NewHandler(services *services.Service) *Handler {
	return &Handler{
		services: services,
	}
}

func (h *Handler) InitRoutes() *gin.Engine {
	router := gin.New()
	docs.SwaggerInfo.BasePath = "/api/v1"
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	wallet := router.Group("/api/v1")
	{
		wallet.POST("/create-wallet", h.CreateWallet)
		wallet.POST("/wallet", h.ChangeBalance)
		wallet.GET("/wallets/:walletId", h.GetBalance)
	}
	return router
}
