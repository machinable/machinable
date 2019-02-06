package models

import (
	"encoding/json"
	"time"

	"github.com/mongodb/mongo-go-driver/bson"
	"github.com/mongodb/mongo-go-driver/bson/objectid"
)

// ProjectAPIKey is a static key used to access resources and collections of the project
type ProjectAPIKey struct {
	ID          objectid.ObjectID `json:"id" bson:"_id,omitempty"`
	KeyHash     string            `json:"-" bson:"key_hash"`
	Created     time.Time         `json:"created" bson:"created"`
	Description string            `json:"description" bson:"description"`
	Read        bool              `json:"read"`
	Write       bool              `json:"write"`
	Role        string            `json:"role"`
}

// MarshalJSON is the custom marshaller for api key structs
func (t ProjectAPIKey) MarshalJSON() ([]byte, error) {
	key := struct {
		ID          string    `json:"id"`
		Description string    `json:"description"`
		Created     time.Time `json:"created"`
		Read        bool      `json:"read"`
		Write       bool      `json:"write"`
		Role        string    `json:"role"`
	}{}

	key.ID = t.ID.Hex()
	key.Description = t.Description
	key.Created = t.Created
	key.Read = t.Read
	key.Write = t.Write
	key.Role = t.Role

	return json.Marshal(&key)
}

// UnmarshalBSON is the custom unmarshaler
func (t *ProjectAPIKey) UnmarshalBSON(bytes []byte) error {
	doc, err := bson.ReadDocument(bytes)
	if err != nil {
		return err
	}

	t.ID = doc.Lookup("_id").ObjectID()
	t.Description = doc.Lookup("description").StringValue()
	t.KeyHash = doc.Lookup("key_hash").StringValue()
	t.Read = doc.Lookup("read").Boolean()
	t.Write = doc.Lookup("write").Boolean()
	t.Role = doc.Lookup("role").StringValue()

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
	Read         bool              `json:"read"`
	Write        bool              `json:"write"`
	Role         string            `json:"role"`
}

// MarshalJSON is the custom marshaller for user structs
func (u ProjectUser) MarshalJSON() ([]byte, error) {
	user := struct {
		ID       string    `json:"id"`
		Username string    `json:"username"`
		Created  time.Time `json:"created"`
		Read     bool      `json:"read"`
		Write    bool      `json:"write"`
		Role     string    `json:"role"`
	}{}

	// Marshal ID to string
	user.ID = u.ID.Hex()
	user.Username = u.Username
	user.Created = u.Created
	user.Read = u.Read
	user.Write = u.Write
	user.Role = u.Role

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
	u.Read = doc.Lookup("read").Boolean()
	u.Write = doc.Lookup("write").Boolean()
	u.Role = doc.Lookup("role").StringValue()

	// This is the only reason we have this Unmarshaler. The default unmarshal is trying to
	// set this as an int64, rather than time.Time
	u.Created = doc.Lookup("created").Time()

	return nil
}
