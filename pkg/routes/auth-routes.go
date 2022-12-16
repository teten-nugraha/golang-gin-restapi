package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/teten-nugraha/golang-gin-project/pkg/handler"
)

func AuthRoutes(route *gin.Engine, authHandler handler.AuthHandler) {
	authRoutes := route.Group("api/auth")
	{
		authRoutes.POST("/login", authHandler.Login)
		authRoutes.POST("/register", authHandler.Register)
	}
}
