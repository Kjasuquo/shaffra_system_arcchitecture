package controller

import (
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"

	"shaffra_assessment/config"
	"shaffra_assessment/internal/pkg/web"
	"shaffra_assessment/internal/service"
)

type Handler struct {
	Config          config.Config
	UserService     service.UserService
	Wg              *sync.WaitGroup
	ReqDurationChan chan string
}

// trackRequestTime serves as a Middleware to track duration in each handler request
func (h *Handler) trackRequestTime(c *gin.Context, wg *sync.WaitGroup, start time.Time) {
	defer wg.Done()
	duration := time.Since(start)
	h.ReqDurationChan <- c.FullPath() + " duration: " + duration.String()
}

// Logger continuously read from the reqDurationChan and logs the request duration
func (h *Handler) Logger() {
	go func() {
		for msg := range h.ReqDurationChan {
			log.Println(msg)
		}
	}()
}

// Ping is used to test the service
func (h *Handler) Ping(c *gin.Context) {
	start := time.Now()
	h.Wg.Add(1)

	defer h.trackRequestTime(c, h.Wg, start)

	web.JSON(c, "pong", http.StatusOK, nil, nil)

}
