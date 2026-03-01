package normalize

import (
	"net"
	"net/url"
	"path"
	"sort"
	"strings"
)

func NormalizeURL(raw string) (string, error) {
	u, err := url.Parse(raw)
	if err != nil {
		return "", err
	}

	u.Scheme = strings.ToLower(u.Scheme)
	u.Host = strings.ToLower(u.Host)

	u.Fragment = ""

	// remove default port
	host, port, err := net.SplitHostPort(u.Host)
	if err == nil {
		if u.Scheme == "http" && port == "80" ||
			u.Scheme == "https" && port == "443" {
			u.Host = host
		}
	}

	// clean path
	u.Path = path.Clean(u.Path)

	if u.Path == "." {
		u.Path = "/"
	}

	// remove tracking query
	query := u.Query()
	for key := range query {
		if strings.HasPrefix(strings.ToLower(key), "utm_") {
			query.Del(key)
		}
	}

	// sort query params
	keys := make([]string, 0, len(query))
	for k := range query {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	sortedQuery := url.Values{}
	for _, k := range keys {
		for _, v := range query[k] {
			sortedQuery.Add(k, v)
		}
	}

	u.RawQuery = sortedQuery.Encode()

	return u.String(), nil
}
