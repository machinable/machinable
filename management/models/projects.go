package models

import (
	"encoding/json"
	"errors"
	"time"

	"github.com/mongodb/mongo-go-driver/bson"
	"github.com/mongodb/mongo-go-driver/bson/objectid"
	"github.com/mongodb/mongo-go-driver/mongo"
)

// reservedProjectSlugs is a list of project slugs that are not allowed to be used, partially so
// we can have reserved sub domains, also so don't muck up our db schema in any way.
var reservedProjectSlugs = map[string]bool{
	"management": true,
	"manage":     true,
	"users":      true,
	"projects":   true,
	"sessions":   true,
	"machinable": true,
}

// ProjectBody is used to unmarshal the JSON body of an incoming request
type ProjectBody struct {
	UserID      string `json:"user_id"`
	Slug        string `json:"slug"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Icon        string `json:"icon"`
	Authn       bool   `json:"authn"`
}

// Validate checks the project body for invalid fields
func (pb *ProjectBody) Validate() error {
	if pb.UserID == "" || pb.Slug == "" || pb.Name == "" || pb.Icon == "" {
		return errors.New("invalid project parameters")
	}
	return nil
}

// DuplicateSlug checks the database for the ProjectBody's slug
func (pb *ProjectBody) DuplicateSlug(col *mongo.Collection) bool {
	// check if the slug is in the `reservedProjectNames`
	if _, ok := reservedProjectSlugs[pb.Slug]; ok {
		// mark as duplicate if this slug is not allowed
		return true
	}
	// look up the user
	documentResult := col.FindOne(
		nil,
		bson.NewDocument(
			bson.EC.String("slug", pb.Slug),
		),
		nil,
	)

	project := make(map[string]interface{})
	// decode project document
	err := documentResult.Decode(project)
	if err != nil {
		// no documents in result, project slug does not already exist
		return false
	}

	// slug already exists
	return true
}

// Project is an application project created and managed by a `User`
type Project struct {
	UserID      objectid.ObjectID `json:"user_id" bson:"user_id"`
	Slug        string            `json:"slug" bson:"slug"`
	Name        string            `json:"name" bson:"name"`
	Description string            `json:"description" bson:"description"`
	Icon        string            `json:"icon" bson:"icon"`
	Created     time.Time         `json:"created" bson:"created"`
	Authn       bool              `json:"authn" bson:"authn"`
}

// MarshalJSON is the custom marshaller for user structs
func (p Project) MarshalJSON() ([]byte, error) {
	project := struct {
		UserID      string    `json:"user_id"`
		Slug        string    `json:"slug"`
		Name        string    `json:"name"`
		Description string    `json:"description"`
		Icon        string    `json:"icon"`
		Created     time.Time `json:"created"`
		Authn       bool      `json:"authn"`
	}{}

	// Marshal ID to string
	project.UserID = p.UserID.Hex()
	project.Slug = p.Slug
	project.Name = p.Name
	project.Description = p.Description
	project.Icon = p.Icon
	project.Created = p.Created
	project.Authn = p.Authn

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
	p.Authn = doc.Lookup("authn").Boolean()

	// This is the only reason we have this Unmarshaler. The default unmarshal is trying to
	// set this as an int64, rather than time.Time
	p.Created = doc.Lookup("created").Time()

	return nil
}
