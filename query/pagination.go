package query

import (
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strconv"
)

const (
	// Limit is the default pagination limit
	Limit = "10"
	// MaxLimit is the maximum allowed page size
	MaxLimit = 100
)

var (
	// Scheme is used to create the links... this should be in a config file
	Scheme = os.Getenv("SCHEME")
)

// Links are the pagination links to a response
type Links struct {
	Self string `json:"self,omitempty"`
	Next string `json:"next,omitempty"`
	Prev string `json:"prev,omitempty"`
}

// GetOffset retrieves the `_offset` query parameter and parses it
func GetOffset(values *url.Values) (int64, error) {
	// get pagination parameters
	offset := values.Get("_offset")

	if offset == "" {
		offset = "0"
	}

	io, err := strconv.Atoi(offset)
	if err != nil {
		return 0, errors.New("invalid offset")
	}
	iOffset := int64(io)

	return iOffset, nil
}

// GetLimit retrieves the `_limit` query parameter
func GetLimit(values *url.Values) (int64, error) {
	limit := values.Get("_limit")

	// Set defaults if necessary
	if limit == "" {
		limit = Limit
	}

	// Parse and validate pagination
	il, err := strconv.Atoi(limit)
	if err != nil || il > MaxLimit || il <= 0 {
		return 0, errors.New("invalid limit")
	}
	iLimit := int64(il)

	return iLimit, nil
}

// NewLinks creates the pagination links
func NewLinks(r *http.Request, limit, offset, max int64) *Links {
	if Scheme == "" {
		Scheme = "http"
	}

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
