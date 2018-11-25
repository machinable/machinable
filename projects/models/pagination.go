package models

import (
	"fmt"
	"net/http"
	"net/url"
)

const (
	Limit    = "10"
	MaxLimit = 100
	Scheme   = "http"
)

// NewLinks creates the pagination links
func NewLinks(r *http.Request, limit, offset, max int64) *Links {
	fqdn := Scheme + "://" + r.Host + r.RequestURI
	req, _ := http.NewRequest("GET", fqdn, nil)

	q := req.URL.Query()
	q.Set("_limit", "%d")
	q.Set("_offset", "%d")
	req.URL.RawQuery, _ = url.QueryUnescape(q.Encode())

	fullURL := req.URL.String()

	selfLink := fmt.Sprintf(fullURL, limit, offset)
	nextLink := ""
	prevLink := ""

	pageMax := (max % limit) + max
	if (limit+offset) < pageMax && limit < max {
		nextLink = fmt.Sprintf(fullURL, limit, offset+limit)
	}

	if offset > 0 {
		prevLink = fmt.Sprintf(fullURL, limit, offset-limit)
	}

	links := &Links{
		Self: selfLink,
		Next: nextLink,
		Prev: prevLink,
	}
	// links := &Links{
	// 	Self: fullURL,
	// }

	return links
}

// Links are the pagination links to a response
type Links struct {
	Self string `json:"self,omitempty"`
	Next string `json:"next,omitempty"`
	Prev string `json:"prev,omitempty"`
}
