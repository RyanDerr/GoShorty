package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type ResponseOKWithDataModel struct {
	Code    int    `json:"code"`
	Data    any    `json:"data"`
	Message string `json:"message"`
}

type ResponseOKModel struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type ResponseErrorModel struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type ResponseErrorCustomModel struct {
	Code    int `json:"code"`
	Message any `json:"message"`
}

func ResponseOKWithData(ctx *gin.Context, data any) {
	response := ResponseOKWithDataModel{
		Code:    1000,
		Data:    data,
		Message: "OK",
	}

	ctx.IndentedJSON(http.StatusOK, response)
}

func ResponseOK(c *gin.Context, message string) {
	response := ResponseOKModel{
		Code:    1000,
		Message: message,
	}

	c.IndentedJSON(http.StatusOK, response)
}

func ResponseCreated(ctx *gin.Context, message string) {
	response := ResponseOKModel{
		Code:    1000,
		Message: message,
	}

	ctx.IndentedJSON(http.StatusCreated, response)
}

func ResponseError(c *gin.Context, err string, code int) {
	response := ResponseErrorModel{
		Code:    99,
		Message: err,
	}

	c.IndentedJSON(code, response)
}

func ResponseCustomError(c *gin.Context, err any, code int) {
	response := ResponseErrorCustomModel{
		Code:    99,
		Message: err,
	}

	c.IndentedJSON(code, response)
}

func ResponseRedirect(c *gin.Context, url string) {
	c.IndentedJSON(http.StatusMovedPermanently, url)
}
