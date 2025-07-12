package app

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/azoma13/archiving-service/config"
	v1 "github.com/azoma13/archiving-service/internal/controller/http/v1"
	"github.com/azoma13/archiving-service/internal/service"
	"github.com/azoma13/archiving-service/pkg/httpserver"
	"github.com/labstack/echo/v4"
)

func Run() {
	err := config.NewConfig()
	if err != nil {
		log.Fatalf("app - Run - config.NewConfig: %v", err)
	}

	log.Println("Initializing services...")
	service := service.NewServices()
	err = os.MkdirAll("archives/", 0777)
	if err != nil {
		log.Fatalf("not create dir with named %s: %v", "archives", err)
	}

	log.Println("Initializing handlers and routes...")
	handler := echo.New()

	v1.NewRouter(handler, service)

	log.Println("Start http server...")
	log.Printf("Server port: %s", config.Cfg.HTTP.Port)
	httpServer := httpserver.New(handler, httpserver.Port(config.Cfg.HTTP.Port))

	log.Println("Configuring graceful shutdown...")
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)
	select {
	case s := <-interrupt:
		log.Println("app - Run - signal: " + s.String())
	case err = <-httpServer.Notify():
		log.Println(fmt.Errorf("app - Run - httpServer.Notify: %w", err))
	}

	log.Println("Shutting down...")
	err = os.RemoveAll("./content/")
	if err != nil {
		log.Println(fmt.Errorf("app - Run - os.RemoveAll: %w", err))
	}
	err = httpServer.Shutdown()
	if err != nil {
		log.Println(fmt.Errorf("app - Run - httpServer.Shutdown: %w", err))
	}
}
