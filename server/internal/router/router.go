package router

import (
	"github.com/gin-gonic/gin"
	"github.com/victoryus84/gorders/internal/handler"
	"github.com/victoryus84/gorders/internal/middleware"
)

// SetupRoutes configures all API routes
func SetupRoutes(router *gin.Engine, hdl *handler.Handlers) {
	
	// --- Rute Core (Sistem) ---
	router.GET("/health", hdl.Core.Check)
	router.GET("/version", hdl.Core.Version)
	router.GET("/ping", hdl.Core.Ping)

	// --- Rute Publice (Auth) ---
	router.POST("/signup", hdl.User.Signup)
	router.POST("/login", hdl.User.Login)

	// --- Rute Protejate v1 ---
	api := router.Group("/api/v1")
	api.Use(middleware.AuthJWT())
	{
		clients := api.Group("/clients")
		{
			clients.GET("", hdl.Client.GetClients)
			clients.POST("", hdl.Client.CreateClient)
			clients.GET("/search", hdl.Client.SearchClients)
			//clients.GET("/:id", hdl.Client.GetClientByID)
		}

		contracts := api.Group("/contracts")
		{
			contracts.GET("", hdl.Contract.GetContracts)
			contracts.POST("", hdl.Contract.CreateContract)
			contracts.GET("/search", hdl.Contract.SearchContracts)
		}
	}
}
