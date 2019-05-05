package models

import (
	"encoding/json"
	"time"

	"github.com/mongodb/mongo-go-driver/bson"
	"github.com/mongodb/mongo-go-driver/bson/objectid"
)

// Project is an application project created and managed by a `User`
type Project struct {
	UserID           objectid.ObjectID `json:"user_id" bson:"user_id"`
	Slug             string            `json:"slug" bson:"slug"`
	Name             string            `json:"name" bson:"name"`
	Description      string            `json:"description" bson:"description"`
	Icon             string            `json:"icon" bson:"icon"`
	Created          time.Time         `json:"created" bson:"created"`
	Authn            bool              `json:"authn" bson:"authn"`
	UserRegistration bool              `json:"user_registration" bson:"user_registration"`
}

// MarshalJSON is the custom marshaller for user structs
func (p Project) MarshalJSON() ([]byte, error) {
	project := struct {
		UserID           string    `json:"user_id"`
		Slug             string    `json:"slug"`
		Name             string    `json:"name"`
		Description      string    `json:"description"`
		Icon             string    `json:"icon"`
		Created          time.Time `json:"created"`
		Authn            bool      `json:"authn"`
		UserRegistration bool      `json:"user_registration"`
	}{}

	// Marshal ID to string
	project.UserID = p.UserID.Hex()
	project.Slug = p.Slug
	project.Name = p.Name
	project.Description = p.Description
	project.Icon = p.Icon
	project.Created = p.Created
	project.Authn = p.Authn
	project.UserRegistration = p.UserRegistration

	return json.Marshal(&project)
}

// UnmarshalBSON is a custom unmarshal function to get the `time.Time` value
func (p *Project) UnmarshalBSON(bytes []byte) error {
	doc, err := bson.ReadDocument(bytes)
	if err != nil {
		return err
	}

	p.UserID = doc.Lookup("user_id").ObjectID()
	p.Slug = doc.Lookup("slug").StringValue()
	p.Name = doc.Lookup("name").StringValue()
	p.Description = doc.Lookup("description").StringValue()
	p.Icon = doc.Lookup("icon").StringValue()
	p.Authn, _ = doc.Lookup("authn").BooleanOK()
	p.UserRegistration, _ = doc.Lookup("user_registration").BooleanOK()

	// This is the only reason we have this Unmarshaler. The default unmarshal is trying to
	// set this as an int64, rather than time.Time
	p.Created = doc.Lookup("created").Time()

	return nil
}
