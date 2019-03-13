package models

// Stats is a struct that contains the size and count of objects in a Resource or Collection
type Stats struct {
	Size  int64 `json:"size"`
	Count int64 `json:"count"`
}

// ResponseTiming is a single request
type ResponseTiming struct {
	Timestamp    int64 `json:"timestamp"`     // timestamp in unix time, i.e. number of seconds elapsed since January 1, 1970 UTC
	ResponseTime int64 `json:"response_time"` // milliseconds
}

// ResponseTimes records the response times of requests to collections and api resources over a 5 minute interval
type ResponseTimes struct {
	Timestamp     int64            `json:"timestamp"`      // timestamp in unix time, i.e. number of seconds elapsed since January 1, 1970 UTC
	ResponseTimes []ResponseTiming `json:"response_times"` // milliseconds
}

// StatusCode records status codes of requests to collections and api resources over a 5 minute interval
type StatusCode struct {
	Timestamp int64            `json:"timestamp"` // timestamp in unix time, i.e. number of seconds elapsed since January 1, 1970 UTC
	Codes     map[string]int64 `json:"codes"`     // a map of status codes to the count
}
