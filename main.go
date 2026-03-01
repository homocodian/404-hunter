package main

import (
	"context"
	"log"
	"time"

	"github.com/homocodian/404-hunter/internal/config"
	"github.com/homocodian/404-hunter/internal/server"
)

func main() {
	config := config.NewConfig(10, 5000)

	newServer := server.CreateServer(config)

	if err := server.RunServer(context.Background(), newServer, 3*time.Second); err != nil {
		log.Fatalf("Server error %v", err)
	}
}
