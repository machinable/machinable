package sessions

import (
	"encoding/json"
	"time"

	"github.com/mongodb/mongo-go-driver/bson"
	"github.com/mongodb/mongo-go-driver/bson/objectid"
)

// Session is a user session model for either the mgmt application or a project
type Session struct {
	ID           objectid.ObjectID `json:"id" bson:"_id,omitempty"`
	UserID       string            `json:"user_id" bson:"user_id"`
	Location     string            `json:"location" bson:"location"`
	Mobile       bool              `json:"mobile" bson:"mobile"`
	IP           string            `json:"ip" bson:"ip"`
	LastAccessed time.Time         `json:"last_accessed" bson:"last_accessed"`
	Browser      string            `json:"browser" bson:"browser"`
	OS           string            `json:"os" bson:"os"`
}

// MarshalJSON customer marshaler for sessions
func (s Session) MarshalJSON() ([]byte, error) {
	session := struct {
		ID           string    `json:"id"`
		UserID       string    `json:"user_id"`
		Location     string    `json:"location"`
		Mobile       bool      `json:"mobile"`
		IP           string    `json:"ip"`
		LastAccessed time.Time `json:"last_accessed"`
		Browser      string    `json:"browser"`
		OS           string    `json:"os"`
	}{}

	session.ID = s.ID.Hex()
	session.UserID = s.UserID
	session.Location = s.Location
	session.Mobile = s.Mobile
	session.IP = s.IP
	session.LastAccessed = s.LastAccessed
	session.Browser = s.Browser
	session.OS = s.OS

	return json.Marshal(&session)
}

// UnmarshalBSON is the custom unmarshaler for sessions
func (s *Session) UnmarshalBSON(bytes []byte) error {
	doc, err := bson.ReadDocument(bytes)
	if err != nil {
		return err
	}

	s.ID = doc.Lookup("_id").ObjectID()
	s.UserID = doc.Lookup("user_id").StringValue()
	s.Location = doc.Lookup("location").StringValue()
	s.Mobile = doc.Lookup("mobile").Boolean()
	s.IP = doc.Lookup("ip").StringValue()
	s.Browser = doc.Lookup("browser").StringValue()
	s.OS = doc.Lookup("os").StringValue()

	// This is the only reason we have this Unmarshaler. The default unmarshal is trying to
	// set this as an int64, rather than time.Time
	s.LastAccessed = doc.Lookup("last_accessed").Time()

	return nil
}
