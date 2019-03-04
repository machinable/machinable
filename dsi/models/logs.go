package models

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
