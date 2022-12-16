package main

import (
	"github.com/gin-gonic/gin"
	"github.com/teten-nugraha/golang-gin-project/pkg/config"
	"github.com/teten-nugraha/golang-gin-project/pkg/handler"
	"github.com/teten-nugraha/golang-gin-project/pkg/repository"
	"github.com/teten-nugraha/golang-gin-project/pkg/routes"
	"github.com/teten-nugraha/golang-gin-project/pkg/service"
	"gorm.io/gorm"
)

// Injection
var (
	db *gorm.DB = config.SetupDBConnection()

	userRepository repository.UserRepository = repository.NewUserRepository(db)
	bookRepository repository.BookRepository = repository.NewBookRepository(db)

	jwtService  service.JWTService  = service.NewJWTService()
	userService service.UserService = service.NewUserService(userRepository)
	bookService service.BookService = service.NewBookService(bookRepository)
	authService service.AuthService = service.NewAuthService(userRepository)

	authHandler handler.AuthHandler = handler.NewAuthHandler(authService, jwtService)
	userHandler handler.UserHandler = handler.NewUserHandler(userService, jwtService)
	bookhandler handler.BookHandler = handler.NewBookHandler(bookService, jwtService)
)

func main() {

	defer config.CloseDBConnection(db)
	r := gin.Default()

	routes.AuthRoutes(r, authHandler)
	routes.UserRoutes(r, userHandler, jwtService)
	routes.BookRoutes(r, bookhandler, jwtService)

	r.Run(":8087")
}
