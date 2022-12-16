package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/teten-nugraha/golang-gin-project/pkg/dto"
	"github.com/teten-nugraha/golang-gin-project/pkg/entity"
	"github.com/teten-nugraha/golang-gin-project/pkg/helpers"
	"github.com/teten-nugraha/golang-gin-project/pkg/service"
	"net/http"
	"strconv"
)

// AuthController interface is a contract what this controller can do
type AuthHandler interface {
	Login(ctx *gin.Context)
	Register(ctxx *gin.Context)
}

type authhandler struct {
	authService service.AuthService
	jwtService  service.JWTService
}

func (c *authhandler) Login(ctx *gin.Context) {
	var loginDTO dto.LoginDTO
	errDTO := ctx.ShouldBind(&loginDTO)
	if errDTO != nil {
		response := helpers.BuildErrorResponse("Failed to process request", errDTO.Error(), helpers.EmptyObj{})
		ctx.AbortWithStatusJSON(http.StatusBadRequest, response)
		return
	}

	authResult := c.authService.VerifyCredential(loginDTO.Email, loginDTO.Password)
	if user, ok := authResult.(entity.User); ok {
		generatedToken := c.jwtService.GenerateToken(strconv.FormatUint(user.ID, 10))
		user.Token = generatedToken
		response := helpers.BuildResponse(true, "OK!", user)
		ctx.JSON(http.StatusOK, response)
		return
	}
}

func (c *authhandler) Register(ctx *gin.Context) {
	var registerDTO dto.RegisterDTO
	errDTO := ctx.ShouldBind(&registerDTO)
	if errDTO != nil {
		response := helpers.BuildErrorResponse("Failed to process request", errDTO.Error(), helpers.EmptyObj{})
		ctx.AbortWithStatusJSON(http.StatusBadRequest, response)
		return
	}

	if !c.authService.IsDuplicateEmail(registerDTO.Email) {
		response := helpers.BuildErrorResponse("Failed to process request", "Duplicate Email", helpers.EmptyObj{})
		ctx.JSON(http.StatusConflict, response)
	} else {
		createdUser := c.authService.CreateUser(registerDTO)
		token := c.jwtService.GenerateToken(strconv.FormatUint(createdUser.ID, 10))
		createdUser.Token = token
		response := helpers.BuildResponse(true, "OK!", createdUser)
		ctx.JSON(http.StatusCreated, response)
	}
}

// NewAuthController creates a new instance of AuthController
func NewAuthHandler(authService service.AuthService, jwtService service.JWTService) AuthHandler {
	return &authhandler{
		authService: authService,
		jwtService:  jwtService,
	}
}
