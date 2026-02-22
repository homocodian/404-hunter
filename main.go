package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"time"

	"golang.org/x/net/html"
)

type RequestUrl struct {
	Url string `json:"url"`
}

var httpClient = &http.Client{
	Timeout: 3 * time.Second,
	Transport: &http.Transport{
		ResponseHeaderTimeout: 2 * time.Second,
		TLSHandshakeTimeout:   2 * time.Second,
	},
}

func main() {
	http.Handle("/", http.FileServer(http.Dir("./static")))
	http.HandleFunc("POST /check", func(w http.ResponseWriter, r *http.Request) {
		var requestUrl RequestUrl
		err := json.NewDecoder(r.Body).Decode(&requestUrl)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		fmt.Println("URL to check: ", requestUrl)

		if !isValidURL(requestUrl.Url) {
			http.Error(w, "Invalid URL", http.StatusBadRequest)
			return
		}

		resp, err := http.Get(requestUrl.Url)
		if err != nil {
			log.Println(err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer resp.Body.Close()

		// Parse HTML content using html.Parse
		html, err := html.Parse(resp.Body)
		if err != nil {
			log.Println(err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		links, deadlinks := getDeadLinks(html)
		fmt.Println("All : ", links)
		fmt.Println("Dead : ", deadlinks)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(200)
		json.NewEncoder(w).Encode(map[string][]string{
			"deadLinks": deadlinks,
		})
	})
	fmt.Println("Listening to port 5000")
	http.ListenAndServe(":5000", nil)

	// resp, err := http.Get()
}

func getDeadLinks(root *html.Node) ([]string, []string) {
	var deadLinks []string
	var links []string

	stack := []*html.Node{root}

	for len(stack) > 0 {
		n := stack[len(stack)-1]
		stack = stack[:len(stack)-1]

		if n.Type == html.ElementNode && n.Data == "a" {
			for _, attr := range n.Attr {
				if attr.Key == "href" {
					links = append(links, attr.Val)
					if !isValidURL(attr.Val) || !isReachable(attr.Val) {
						deadLinks = append(deadLinks, attr.Val)
					}
				}
			}
		}
		for child := n.FirstChild; child != nil; child = child.NextSibling {
			stack = append(stack, child)
		}
	}

	return links, deadLinks
}

func isValidURL(urlToCheck string) bool {
	u, err := url.ParseRequestURI(urlToCheck)
	if err != nil {
		log.Println(err.Error())
		return false
	}
	fmt.Println("Url: ", u)
	return u.Scheme != "" && u.Host != "" && u.Hostname() != "localhost" && u.Hostname() != "127.0.0.1"
}

func isReachable(urlToCheck string) bool {
	resp, err := httpClient.Head(urlToCheck)
	if err != nil {
		log.Println(err.Error())
		return false
	}

	defer resp.Body.Close()
	fmt.Println(urlToCheck, " : ", resp.StatusCode)
	return resp.StatusCode < 400
}
