package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/teten-nugraha/golang-gin-project/pkg/handler"
	"github.com/teten-nugraha/golang-gin-project/pkg/middleware"
	"github.com/teten-nugraha/golang-gin-project/pkg/service"
)

func UserRoutes(route *gin.Engine, userHandler handler.UserHandler, jwtService service.JWTService) {
	userRoutes := route.Group("api/user", middleware.AuthorizeJWT(jwtService))
	{
		userRoutes.GET("/profile", userHandler.Profile)
		userRoutes.PUT("/profile", userHandler.Update)
	}
}
