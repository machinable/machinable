package models

import (
	"encoding/json"
	"time"

	"github.com/mongodb/mongo-go-driver/bson"
)

// User is a user of the application
type User struct {
	ID           string    `json:"id" bson:"_id,omitempty"`
	Username     string    `json:"username" bson:"username"`
	Email        string    `json:"-" bson:"username"`
	PasswordHash string    `json:"-" bson:"password_hash"`
	Created      time.Time `json:"created" bson:"created"`
}

// MarshalJSON is a custom json marshal function
func (u User) MarshalJSON() ([]byte, error) {
	user := struct {
		ID       string    `json:"id"`
		Username string    `json:"username"`
		Created  time.Time `json:"created"`
	}{}

	// Marshal ID to string
	user.ID = u.ID
	user.Username = u.Username
	user.Created = u.Created

	return json.Marshal(&user)
}

// UnmarshalBSON is a custom unmarshal function to get the `time.Time` value
func (u *User) UnmarshalBSON(bytes []byte) error {
	doc, err := bson.ReadDocument(bytes)
	if err != nil {
		return err
	}

	u.ID = doc.Lookup("_id").ObjectID().Hex()
	u.Username = doc.Lookup("username").StringValue()
	u.PasswordHash = doc.Lookup("password_hash").StringValue()

	// This is the only reason we have this Unmarshaler. The default unmarshal is trying to
	// set this as an int64, rather than time.Time
	u.Created = doc.Lookup("created").Time()

	return nil
}
