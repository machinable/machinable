package models

import "time"

// Log is any user/api key initiated event that should be recorded
type Log struct {
	Event       string    `json:"event"`
	Created     time.Time `json:"created"`
	Initiator   string    `json:"initiator"`
	InitiatorID string    `json:"initiator_id"`
	TargetID    string    `json:"target_id"`
}
