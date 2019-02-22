package models

import (
	"fmt"
	"net/http"
	"net/url"
)

const (
	// Limit is the default pagination limit
	Limit = "10"
	// MaxLimit is the maximum allowed page size
	MaxLimit = 100
	// Scheme is used to create the links... this should be in a config file
	Scheme = "http"
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
		prevOffset := offset - limit
		if prevOffset < 0 {
			prevOffset = 0
		}
		prevLink = fmt.Sprintf(fullURL, limit, prevOffset)
	}

	links := &Links{
		Self: selfLink,
		Next: nextLink,
		Prev: prevLink,
	}

	return links
}

// Links are the pagination links to a response
type Links struct {
	Self string `json:"self,omitempty"`
	Next string `json:"next,omitempty"`
	Prev string `json:"prev,omitempty"`
}
