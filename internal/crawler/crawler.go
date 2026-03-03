package crawler

import (
	"context"
	"log"
	"net/http"
	"net/url"
	"strings"
	"sync"

	"github.com/homocodian/404-hunter/internal/normalize"
	"github.com/homocodian/404-hunter/internal/parser"
	"github.com/homocodian/404-hunter/internal/storage"
	"golang.org/x/net/html"
)

type Crawler struct {
	baseHost string
	client   *http.Client
	workers  int
	queue    chan string
	visited  *storage.InMemoryVisited
	broken   *storage.InMemoryBroken
}

func NewCrawler(hostname string, client *http.Client, worker int) *Crawler {
	return &Crawler{
		baseHost: hostname,
		client:   client,
		workers:  worker,
		queue:    make(chan string, 500),
		visited:  storage.NewInMemoryVisited(),
		broken:   storage.NewInMemoryBroken(),
	}
}

func (c *Crawler) Start(ctx context.Context, startURL string) []string {

	var wg sync.WaitGroup

	// Seed the start URL
	if c.visited.Add(startURL) {
		wg.Add(1)
		go func() {
			select {
			case c.queue <- startURL:
			case <-ctx.Done():
				wg.Done()
			}
		}()
	}

	for range c.workers {
		go func() {
			for {
				select {
				case <-ctx.Done():
					return
				case urlStr, ok := <-c.queue:
					if !ok {
						return
					}
					c.processWithGW(ctx, urlStr, &wg)
					wg.Done()
				}
			}
		}()
	}

	done := make(chan struct{})
	go func() {
		wg.Wait()
		close(c.queue)
		close(done)
	}()

	select {
	case <-ctx.Done():
	case <-done:
	}

	return c.broken.GetAll()
}

func (c *Crawler) processWithGW(ctx context.Context, rawURL string, wg *sync.WaitGroup) {
	baseURL, err := url.Parse(rawURL)
	if err != nil {
		log.Println(err)
		return
	}

	if !c.isAllowed(baseURL.Hostname()) {
		req, _ := http.NewRequestWithContext(ctx, "HEAD", rawURL, nil)
		resp, err := c.client.Do(req)
		if err != nil {
			if ctx.Err() == nil {
				c.broken.Add(rawURL)
			}
			log.Println(err)
			return
		}
		defer resp.Body.Close()

		if resp.StatusCode >= 400 {
			c.broken.Add(rawURL)
		}
		return
	}

	req, err := http.NewRequestWithContext(ctx, "GET", rawURL, nil)
	if err != nil {
		log.Println(err)
		return
	}

	resp, err := c.client.Do(req)
	if err != nil {
		if ctx.Err() == nil {
			log.Println("Request error:", err)
		}
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		c.broken.Add(rawURL)
		return
	}

	contentType := resp.Header.Get("Content-Type")

	if !strings.HasPrefix(contentType, "text/html") {
		return
	}

	html, err := html.Parse(resp.Body)
	if err != nil {
		log.Println(err)
		return
	}

	links := parser.ExtractLinks(html)

	for _, link := range links {
		abs, err := parser.ResolveURL(baseURL, link)
		if err != nil {
			continue
		}
		norm, err := normalize.NormalizeURL(abs.String())
		if err != nil {
			continue
		}

		if c.visited.Add(norm) {
			wg.Add(1)
			go func(l string) {
				select {
				case c.queue <- l:
				case <-ctx.Done():
					wg.Done()
				}
			}(norm)
		}
	}
}

func (c *Crawler) isAllowed(hostname string) bool {
	if hostname == c.baseHost &&
		hostname != "localhost" &&
		hostname != "127.0.0.1" {
		return true
	}
	return false
}
