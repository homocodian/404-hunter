package main

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/homocodian/404-hunter/internal/argparser"
	"github.com/homocodian/404-hunter/internal/cli"
	"github.com/homocodian/404-hunter/internal/config"
	"github.com/homocodian/404-hunter/internal/server"
)

func main() {
	args, err := argparser.Parse()
	if err != nil {
		log.Fatal(err)
	}

	workers := args.GetInt("w", 5)

	if workers <= 0 && workers != -1 {
		log.Fatal("Invalid workers value")
	}

	switch os.Args[1] {
	case "web":
		port := args.GetInt("p", 5000)

		cfg := config.NewServerConfig(workers, port)

		newServer := server.CreateServer(cfg)

		if err := server.RunServer(context.Background(), newServer, 3*time.Second); err != nil {
			log.Fatalf("Server error %v", err)
		}

	default:
		var opts []config.Option

		outputFile, ok := args["o"]
		if ok && outputFile == "" {
			log.Fatal("Invalid filename/location")
		}

		// store output filename in config
		if ok && outputFile != "" {
			opts = append(opts, config.WithOutputFile(outputFile))
		}

		cfg := config.NewConfig(workers, opts...)

		if err := cli.RunCli(cfg, os.Args[1]); err != nil {
			log.Fatal(err)
		}
	}
}
