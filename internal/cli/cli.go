package cli

import (
	"context"
	"errors"
	"fmt"
	"net/url"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/homocodian/404-hunter/internal/client"
	"github.com/homocodian/404-hunter/internal/config"
	"github.com/homocodian/404-hunter/internal/crawler"
)

const (
	fileModeUserRW = 0644
)

func RunCli(config *config.Config, targetURL string) error {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	httpClient := client.NewHttpClient(config)

	parsedURL, err := url.Parse(targetURL)
	if err != nil {
		return errors.New(fmt.Sprint("Failed to parse url:", err.Error()))
	}

	c := crawler.NewCrawler(parsedURL.Hostname(), httpClient, config.Workers)
	brokenLinks := c.Start(ctx, targetURL)

	if ctx.Err() != nil {
		return errors.New("Stopped due to user/system signal")
	}

	if config.OutputFile != "" {
		return os.WriteFile(config.OutputFile, []byte(strings.Join(brokenLinks, ",")), fileModeUserRW)
	}

	fmt.Println(brokenLinks)

	return nil
}
