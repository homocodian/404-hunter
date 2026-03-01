package parser

import "net/url"

func ResolveURL(baseURL *url.URL, relativeURL string) (*url.URL, error) {
	relative, err := url.Parse(relativeURL)
	if err != nil {
		return nil, err
	}
	return baseURL.ResolveReference(relative), nil
}
