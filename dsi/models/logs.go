package models

import "strconv"

// Log is any user/api key initiated event that should be recorded
type Log struct {
	Event         string `json:"event"`
	StatusCode    int    `json:"status_code"`
	Created       int64  `json:"created"`
	Initiator     string `json:"initiator"`
	InitiatorType string `json:"initiator_type"`
	InitiatorID   string `json:"initiator_id"`
	TargetID      string `json:"target_id"`
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
