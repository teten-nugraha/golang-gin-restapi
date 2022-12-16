package handler

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/teten-nugraha/golang-gin-project/pkg/dto"
	"github.com/teten-nugraha/golang-gin-project/pkg/helpers"
	"github.com/teten-nugraha/golang-gin-project/pkg/service"
	"net/http"
	"strconv"
)

type UserHandler interface {
	Update(context *gin.Context)
	Profile(context *gin.Context)
}

type userHandler struct {
	userService service.UserService
	jwtService  service.JWTService
}

func (c *userHandler) Update(context *gin.Context) {
	var updateUserDTO dto.UserUpdateDTO
	errDTO := context.ShouldBind(&updateUserDTO)
	if errDTO != nil {
		res := helpers.BuildErrorResponse("Failed to process request", errDTO.Error(), helpers.EmptyObj{})
		context.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	authHeader := context.GetHeader("Authorization")
	token, errToken := c.jwtService.ValidateToken(authHeader)
	if errToken != nil {
		panic(errToken.Error())
	}
	claims := token.Claims.(jwt.MapClaims)
	id, err := strconv.ParseUint(fmt.Sprintf("%v", claims["user_id"]), 10, 64)
	if err != nil {
		panic(err.Error())
	}

	updateUserDTO.ID = id
	u := c.userService.Update(updateUserDTO)
	res := helpers.BuildResponse(true, "Ok", u)
	context.JSON(http.StatusOK, res)

}

func (c *userHandler) Profile(context *gin.Context) {
	authHeader := context.GetHeader("Authorization")
	token, errToken := c.jwtService.ValidateToken(authHeader)
	if errToken != nil {
		panic(errToken.Error())
	}
	claims := token.Claims.(jwt.MapClaims)
	id := fmt.Sprintf("%v", claims["user_id"])

	user := c.userService.Profile(id)
	res := helpers.BuildResponse(true, "OK", user)
	context.JSON(http.StatusOK, res)

}

// NewUserHandler is creating anew instance of UserHandler
func NewUserHandler(userService service.UserService, jwtService service.JWTService) UserHandler {
	return &userHandler{
		userService: userService,
		jwtService:  jwtService,
	}
}
