package models

import (
	"encoding/json"
	"time"

	"github.com/mongodb/mongo-go-driver/bson"
	"github.com/mongodb/mongo-go-driver/bson/objectid"
)

// ProjectAPIToken is a static token used to access resources and collections of the project
type ProjectAPIToken struct {
	ID          objectid.ObjectID `json:"id" bson:"_id,omitempty"`
	TokenHash   string            `json:"-" bson:"token_hash"`
	Created     time.Time         `json:"created" bson:"created"`
	Description string            `json:"description" bson:"description"`
}

// MarshalJSON is the custom marshaller for api token structs
func (t ProjectAPIToken) MarshalJSON() ([]byte, error) {
	token := struct {
		ID          string    `json:"id"`
		Description string    `json:"description"`
		Created     time.Time `json:"created"`
	}{}

	token.ID = t.ID.Hex()
	token.Description = t.Description
	token.Created = t.Created

	return json.Marshal(&token)
}

// UnmarshalBSON is the custom unmarshaler
func (t *ProjectAPIToken) UnmarshalBSON(bytes []byte) error {
	doc, err := bson.ReadDocument(bytes)
	if err != nil {
		return err
	}

	t.ID = doc.Lookup("_id").ObjectID()
	t.Description = doc.Lookup("description").StringValue()
	t.TokenHash = doc.Lookup("token_hash").StringValue()

	// This is the only reason we have this Unmarshaler. The default unmarshal is trying to
	// set this as an int64, rather than time.Time
	t.Created = doc.Lookup("created").Time()

	return nil
}

// ProjectUser is a user of a project. A user can access resources and collections of the project.
type ProjectUser struct {
	ID           objectid.ObjectID `json:"id" bson:"_id,omitempty"`
	Username     string            `json:"username" bson:"username"`
	PasswordHash string            `json:"-" bson:"password_hash"`
	Created      time.Time         `json:"created" bson:"created"`
}

// MarshalJSON is the custom marshaller for user structs
func (u ProjectUser) MarshalJSON() ([]byte, error) {
	user := struct {
		ID       string    `json:"id"`
		Username string    `json:"username"`
		Created  time.Time `json:"created"`
	}{}

	// Marshal ID to string
	user.ID = u.ID.Hex()
	user.Username = u.Username
	user.Created = u.Created

	return json.Marshal(&user)
}

// UnmarshalBSON is the custom unmarshaler
func (u *ProjectUser) UnmarshalBSON(bytes []byte) error {
	doc, err := bson.ReadDocument(bytes)
	if err != nil {
		return err
	}

	u.ID = doc.Lookup("_id").ObjectID()
	u.Username = doc.Lookup("username").StringValue()
	u.PasswordHash = doc.Lookup("password_hash").StringValue()

	// This is the only reason we have this Unmarshaler. The default unmarshal is trying to
	// set this as an int64, rather than time.Time
	u.Created = doc.Lookup("created").Time()

	return nil
}

// NewProjectUser is the JSON structure of a new user request
type NewProjectUser struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// NewProjectToken is the JSON structure of a new api token request
type NewProjectToken struct {
	Token       string `json:"token"`
	Description string `json:"description"`
}
