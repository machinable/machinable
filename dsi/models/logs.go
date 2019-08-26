package models

import "strconv"

const (
	EndpointResource   string = "resource"
	EndpointCollection string = "collection"
)

// Log is any user/api key initiated event that should be recorded
type Log struct {
	ID             string `json:"id"`
	ProjectID      string `json:"project_id"`
	EndpointType   string `json:"endpoint_type"`
	Verb           string `json:"verb"`
	Path           string `json:"path"`
	StatusCode     int    `json:"status_code"`
	Created        int64  `json:"created"`
	AlignedCreated int64  `json:"aligned"`
	ResponseTime   int64  `json:"response_time"`
	Initiator      string `json:"initiator"`
	InitiatorType  string `json:"initiator_type"`
	InitiatorID    string `json:"initiator_id"`
	TargetID       string `json:"target_id"`
}

// FieldAsTypedInterface returns the field value as an interface with the proper type
func FieldAsTypedInterface(field string, val string) (interface{}, error) {
	if field == "status_code" {
		return strconv.Atoi(val)
	} else if field == "created" {
		return strconv.ParseInt(val, 10, 64)
	}

	return val, nil
}

// IsValidLogField verifies that a field exists for the `Log` object
func IsValidLogField(field string) bool {
	fields := []string{"event", "status_code", "created", "initiator", "initiator_type", "initiator_id", "target_id"}

	for _, f := range fields {
		if f == field {
			return true
		}
	}

	return false
}
