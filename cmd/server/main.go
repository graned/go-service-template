package main

import (
	"github.com/graned/go-service-template/internal/transport/rest"
	"github.com/graned/go-service-template/internal/user"
	"log"
)

func main() {
	userService := user.NewService()

	restServer := rest.New(":3000", userService)

	if err := restServer.Run(); err != nil {
		log.Fatal(err)
	}
}
