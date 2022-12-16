package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/teten-nugraha/golang-gin-project/pkg/handler"
	"github.com/teten-nugraha/golang-gin-project/pkg/middleware"
	"github.com/teten-nugraha/golang-gin-project/pkg/service"
)

func BookRoutes(route *gin.Engine, bookHandler handler.BookHandler, jwtService service.JWTService) {
	bookRoutes := route.Group("api/books", middleware.AuthorizeJWT(jwtService))
	{
		bookRoutes.GET("/", bookHandler.All)
		bookRoutes.POST("/", bookHandler.Insert)
		bookRoutes.GET("/:id", bookHandler.FindByID)
		bookRoutes.PUT("/:id", bookHandler.Update)
		bookRoutes.DELETE("/:id", bookHandler.Delete)
	}
}
