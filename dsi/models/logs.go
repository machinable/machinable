package models

import (
	"time"

	"github.com/mongodb/mongo-go-driver/bson"
)

// Log is any user/api key initiated event that should be recorded
type Log struct {
	Event         string    `json:"event"`
	StatusCode    int       `json:"status_code"`
	Created       time.Time `json:"created"`
	Initiator     string    `json:"initiator"`
	InitiatorType string    `json:"initiator_type"`
	InitiatorID   string    `json:"initiator_id"`
	TargetID      string    `json:"target_id"`
}

// UnmarshalBSON is the custom unmarshaler
func (l *Log) UnmarshalBSON(bytes []byte) error {
	doc, err := bson.ReadDocument(bytes)
	if err != nil {
		return err
	}

	l.Event = doc.Lookup("event").StringValue()
	l.StatusCode = int(doc.Lookup("statuscode").Int64())
	l.Initiator = doc.Lookup("initiator").StringValue()
	l.InitiatorType = doc.Lookup("initiator_type").StringValue()
	l.InitiatorID = doc.Lookup("initiatorid").StringValue()
	l.TargetID = doc.Lookup("targetid").StringValue()

	// This is the only reason we have this Unmarshaler. The default unmarshal is trying to
	// set this as an int64, rather than time.Time
	l.Created = doc.Lookup("created").Time()

	return nil
}
