package web

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func JSON(c *gin.Context, message string, status int, data interface{}, err error) {
	errMessage := ""
	if err != nil {
		errMessage = err.Error()
	}
	responsedata := gin.H{
		"message": message,
		"data":    data,
		"errors":  errMessage,
		"status":  http.StatusText(status),
	}

	c.JSON(status, responsedata)
}

type Response struct {
	Data interface{}
	Err  error
}
