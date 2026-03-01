package parser

import (
	"net/url"

	"golang.org/x/net/html"
)

func ExtractLinks(root *html.Node) []string {
	var links []string

	stack := []*html.Node{root}

	for len(stack) > 0 {
		n := stack[len(stack)-1]
		stack = stack[:len(stack)-1]

		if n.Type == html.ElementNode && n.Data == "a" {
			for _, attr := range n.Attr {
				if attr.Key == "href" {
					links = append(links, attr.Val)
				}
			}
		}
		for child := n.FirstChild; child != nil; child = child.NextSibling {
			stack = append(stack, child)
		}
	}

	return links
}

func ParseURL(raw string) (*url.URL, error) {
	return url.Parse(raw)
}
