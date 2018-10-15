package models

import (
	"encoding/json"
	"errors"
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
	Read        bool              `json:"read"`
	Write       bool              `json:"write"`
}

// MarshalJSON is the custom marshaller for api token structs
func (t ProjectAPIToken) MarshalJSON() ([]byte, error) {
	token := struct {
		ID          string    `json:"id"`
		Description string    `json:"description"`
		Created     time.Time `json:"created"`
		Read        bool      `json:"read"`
		Write       bool      `json:"write"`
	}{}

	token.ID = t.ID.Hex()
	token.Description = t.Description
	token.Created = t.Created
	token.Read = t.Read
	token.Write = t.Write

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
	t.Read = doc.Lookup("read").Boolean()
	t.Write = doc.Lookup("write").Boolean()

	// This is the only reason we have this Unmarshaler. The default unmarshal is trying to
	// set this as an int64, rather than time.Time
	t.Created = doc.Lookup("created").Time()

	return nil
}

// ProjectSession is a user session for a project
type ProjectSession struct {
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
func (s ProjectSession) MarshalJSON() ([]byte, error) {
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
func (s *ProjectSession) UnmarshalBSON(bytes []byte) error {
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
	s.LastAccessed = doc.Lookup("created").Time()

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
}

// MarshalJSON is the custom marshaller for user structs
func (u ProjectUser) MarshalJSON() ([]byte, error) {
	user := struct {
		ID       string    `json:"id"`
		Username string    `json:"username"`
		Created  time.Time `json:"created"`
		Read     bool      `json:"read"`
		Write    bool      `json:"write"`
	}{}

	// Marshal ID to string
	user.ID = u.ID.Hex()
	user.Username = u.Username
	user.Created = u.Created
	user.Read = u.Read
	user.Write = u.Write

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

	// This is the only reason we have this Unmarshaler. The default unmarshal is trying to
	// set this as an int64, rather than time.Time
	u.Created = doc.Lookup("created").Time()

	return nil
}

// NewProjectUser is the JSON structure of a new user request
type NewProjectUser struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Read     bool   `json:"read"`
	Write    bool   `json:"write"`
}

// Validate checks that the new user has a username and password.
func (u *NewProjectUser) Validate() error {
	if u.Username == "" || u.Password == "" {
		return errors.New("invalid username or password")
	}
	return nil
}

// NewProjectToken is the JSON structure of a new api token request
type NewProjectToken struct {
	Token       string `json:"token"`
	Description string `json:"description"`
	Read        bool   `json:"read"`
	Write       bool   `json:"write"`
}

// Validate checks that the new token is not empty
func (u *NewProjectToken) Validate() error {
	if u.Token == "" {
		return errors.New("invalid token")
	}
	return nil
}
