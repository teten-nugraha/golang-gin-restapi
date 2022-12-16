package handler

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/teten-nugraha/golang-gin-project/pkg/dto"
	"github.com/teten-nugraha/golang-gin-project/pkg/entity"
	"github.com/teten-nugraha/golang-gin-project/pkg/helpers"
	"github.com/teten-nugraha/golang-gin-project/pkg/service"
	"net/http"
	"strconv"
)

// BookController is a ...
type BookHandler interface {
	All(context *gin.Context)
	FindByID(context *gin.Context)
	Insert(context *gin.Context)
	Update(context *gin.Context)
	Delete(context *gin.Context)
}

type bookHandler struct {
	bookService service.BookService
	jwtService  service.JWTService
}

func (b *bookHandler) All(context *gin.Context) {
	var books []entity.Book = b.bookService.All()
	res := helpers.BuildResponse(true, "OK", books)
	context.JSON(http.StatusOK, res)
}

func (b *bookHandler) FindByID(context *gin.Context) {
	id, err := strconv.ParseUint(context.Param("id"), 0, 0)
	if err != nil {
		res := helpers.BuildErrorResponse("No param id was found", err.Error(), helpers.EmptyObj{})
		context.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	var book entity.Book = b.bookService.FindByID(id)
	if (book == entity.Book{}) {
		res := helpers.BuildErrorResponse("Data not found", "No data with given id ", helpers.EmptyObj{})
		context.AbortWithStatusJSON(http.StatusNotFound, res)
	} else {
		res := helpers.BuildResponse(true, "OK", book)
		context.JSON(http.StatusOK, res)
	}
}

func (b *bookHandler) Insert(context *gin.Context) {
	var bookCreateDTO dto.BookCreateDTO
	errDTO := context.ShouldBind(&bookCreateDTO)
	if errDTO != nil {
		res := helpers.BuildErrorResponse("Failed to process request", errDTO.Error(), helpers.EmptyObj{})
		context.JSON(http.StatusBadRequest, res)
	} else {
		authHeader := context.GetHeader("Authorization")
		userID := b.getUserIDByToken(authHeader)
		convertedID, err := strconv.ParseUint(userID, 10, 64)
		if err != nil {
			bookCreateDTO.UserID = convertedID
		}

		result := b.bookService.Insert(bookCreateDTO)
		response := helpers.BuildResponse(true, "OK", result)
		context.JSON(http.StatusCreated, response)
	}
}

func (b *bookHandler) Update(context *gin.Context) {
	var bookUpdateDTO dto.BookUpdateDTO
	errDTO := context.ShouldBind(&bookUpdateDTO)
	if errDTO != nil {
		res := helpers.BuildErrorResponse("Failed to process request", errDTO.Error(), helpers.EmptyObj{})
		context.JSON(http.StatusBadRequest, res)
		return
	}

	authHeader := context.GetHeader("Authorization")
	token, errToken := b.jwtService.ValidateToken(authHeader)
	if errToken != nil {
		panic(errToken.Error())
	}
	claims := token.Claims.(jwt.MapClaims)
	userID := fmt.Sprintf("%v", claims["user_id"])
	if b.bookService.IsAllowedToEdit(userID, bookUpdateDTO.ID) {
		id, errID := strconv.ParseUint(userID, 10, 64)
		if errID == nil {
			bookUpdateDTO.UserID = id
		}
		result := b.bookService.Update(bookUpdateDTO)
		response := helpers.BuildResponse(true, "OK", result)
		context.JSON(http.StatusOK, response)
	} else {
		response := helpers.BuildErrorResponse("You dont have permission", "You are not the owner", helpers.EmptyObj{})
		context.JSON(http.StatusForbidden, response)
	}
}

func (b *bookHandler) Delete(context *gin.Context) {
	var book entity.Book
	id, err := strconv.ParseUint(context.Param("id"), 0, 0)
	if err != nil {
		response := helpers.BuildErrorResponse("failed to get id", "No param id were found", helpers.EmptyObj{})
		context.JSON(http.StatusBadRequest, response)
	}
	book.ID = id

	authHeader := context.GetHeader("Authorization")
	token, errToken := b.jwtService.ValidateToken(authHeader)
	if errToken != nil {
		panic(errToken.Error())
	}
	claims := token.Claims.(jwt.MapClaims)
	userID := fmt.Sprintf("%v", claims["user_id"])

	if b.bookService.IsAllowedToEdit(userID, book.ID) {
		b.bookService.Delete(book)
		response := helpers.BuildResponse(true, "Deleted", helpers.EmptyObj{})
		context.JSON(http.StatusOK, response)
	} else {
		response := helpers.BuildErrorResponse("You dont have permission", "You are not the owner", helpers.EmptyObj{})
		context.JSON(http.StatusForbidden, response)
	}
}

func (b *bookHandler) getUserIDByToken(token string) string {
	aToken, err := b.jwtService.ValidateToken(token)
	if err != nil {
		panic(err.Error())
	}

	claims := aToken.Claims.(jwt.MapClaims)
	id := fmt.Sprintf("%v", claims["user_id"])
	return id
}

// NewBookController create a new instances of BoookController
func NewBookHandler(bookServ service.BookService, jwtServ service.JWTService) BookHandler {
	return &bookHandler{
		bookService: bookServ,
		jwtService:  jwtServ,
	}
}
