package server

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"shaffra_assessment/config"
	"shaffra_assessment/internal/controller"
	"shaffra_assessment/internal/repository"
	"shaffra_assessment/internal/service"
)

func Start() {
	var (
		DB = new(repository.PostgresDB)
	)
	conf, err := config.InitDBConfigs()
	if err != nil {
		log.Fatal(err)
	}
	userService := service.NewService(DB)

	h := &controller.Handler{
		Config:          *conf,
		UserService:     userService,
		Wg:              &sync.WaitGroup{},
		ReqDurationChan: make(chan string),
	}

	err = DB.Connect(*conf)
	if err != nil {
		log.Fatalf("Error trying to connect to user: %v", err)
	}

	r := SetupRouter(h)

	PORT := fmt.Sprintf(":%s", conf.ServicePort)
	if PORT == ":" {
		PORT = ":5053"
	}
	srv := &http.Server{
		Addr:    PORT,
		Handler: r,
	}

	// Initializing the server in a goroutine so that
	// it won't block the graceful shutdown handling below
	go func() {
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	log.Printf("Server started on %s\n", PORT)
	gracefulShutdown(srv)
}

func gracefulShutdown(srv *http.Server) {
	// Wait for interrupt signal to gracefully shutdown the server with
	// a timeout of 5 seconds.
	quit := make(chan os.Signal, 1)
	// kill (no param) default send syscall.SIGTERM
	// kill -2 is syscall.SIGINT
	// kill -9 is syscall.SIGKILL but can't be catch, so don't need add it
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")

	// The context is used to inform the server it has 5 seconds to finish
	// the request it is currently handling
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown:", err)
	}
	log.Println("Server exiting")
}
