package main

import (
	applogger "github.com/graned/go-service-template/internal/logger"
	"github.com/graned/go-service-template/internal/transport/rest"
	"github.com/graned/go-service-template/internal/user"
	"log"
)

func main() {
	logger := applogger.New()

	userService := user.NewService()

	restServer := rest.New(
		":3000",
		logger,
		userService,
	)

	if err := restServer.Run(); err != nil {
		log.Fatal(err)
	}
}
