package server

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/homocodian/404-hunter/internal/config"
	"github.com/homocodian/404-hunter/internal/crawler"
	"github.com/homocodian/404-hunter/internal/parser"
)

type RequestUrl struct {
	Url string `json:"url"`
}

func CreateServer(config *config.Config) *http.Server {
	var httpClient = &http.Client{
		Timeout: 3 * time.Second,
		Transport: &http.Transport{
			ResponseHeaderTimeout: 2 * time.Second,
			TLSHandshakeTimeout:   2 * time.Second,
			MaxIdleConns:          1000,
			MaxIdleConnsPerHost:   config.Workers,
			IdleConnTimeout:       90 * time.Second,
		},
	}

	mux := http.NewServeMux()

	mux.Handle("/", http.FileServer(http.Dir("./static")))
	mux.HandleFunc("POST /check", func(w http.ResponseWriter, r *http.Request) {
		var requestUrl RequestUrl
		err := json.NewDecoder(r.Body).Decode(&requestUrl)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		parsedURL, err := parser.ParseURL(requestUrl.Url)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		startTime := time.Now()

		c := crawler.NewCrawler(parsedURL.Hostname(), httpClient, config.Workers)
		brokenLinks := c.Start(r.Context(), requestUrl.Url)

		endTime := time.Since(startTime)

		switch {
		case endTime > time.Minute:
			fmt.Printf("%d minutes\n", endTime/time.Minute)
		case endTime > time.Second:
			fmt.Printf("%d seconds\n", endTime/time.Second)
		default:
			fmt.Printf("%d milliseconds\n", endTime/time.Millisecond)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(200)
		json.NewEncoder(w).Encode(map[string][]string{
			"deadLinks": brokenLinks,
		})
	})

	return &http.Server{
		Addr:         fmt.Sprint(":", config.Port),
		Handler:      mux,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 30 * time.Minute,
	}
}

func RunServer(ctx context.Context, server *http.Server, shutdownTimeout time.Duration) error {

	serverError := make(chan error, 1)

	go func() {
		fmt.Println("Listening to ", server.Addr)
		if err := server.ListenAndServe(); errors.Is(err, http.ErrServerClosed) {
			serverError <- err
		}
		close(serverError)
	}()

	stop := make(chan os.Signal, 1)

	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	select {
	case err := <-serverError:
		return err
	case <-stop:
		log.Println("Shutdown signal received...")
	case <-ctx.Done():
		log.Println("Contex cancelled")
	}

	shutdownContext, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()

	if err := server.Shutdown(shutdownContext); err != nil {
		if closeError := server.Close(); closeError != nil {
			return errors.Join(err, closeError)
		}
		return err
	}

	log.Println("Server exited gracefully")
	return nil
}
