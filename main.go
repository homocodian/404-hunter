package main

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/homocodian/404-hunter/internal/cli"
	"github.com/homocodian/404-hunter/internal/config"
	"github.com/homocodian/404-hunter/internal/server"
)

func main() {
	if len(os.Args) < 2 || os.Args[1] == "" {
		log.Fatal("Invalid arguments, see --help for more info.")
	}

	config := config.NewConfig(10, 5000)

	switch os.Args[1] {
	case "web":
		newServer := server.CreateServer(config)

		if err := server.RunServer(context.Background(), newServer, 3*time.Second); err != nil {
			log.Fatalf("Server error %v", err)
		}

	default:
		if err := cli.RunCli(config, os.Args[1]); err != nil {
			log.Fatal(err)
		}
	}
}
