package server

import (
	"log"

	"github.com/gin-gonic/gin"

	"shaffra_assessment/internal/controller"
)

func DefineRoutes(handler *controller.Handler) *gin.Engine {
	log.Println("Routes defined")

	router := gin.Default()

	handler.Logger()

	{
		router.GET("/ping", handler.Ping)
		router.POST("/users", handler.CreateUser)
		router.GET("/users/:id", handler.GetUser)
		router.PUT("/users/:id", handler.UpdateUser)
		router.DELETE("/users/:id", handler.DeleteUser)
	}

	return router
}

func SetupRouter(h *controller.Handler) *gin.Engine {
	log.Println("Router setup")
	r := DefineRoutes(h)

	return r
}
