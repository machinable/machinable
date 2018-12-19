package models

import (
	"errors"
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

// ReservedSlug verifies the slug is not reserved. Returns true if reserved, false otherwise
func (pb *ProjectBody) ReservedSlug() bool {
	// check if the slug is in the `reservedProjectNames`
	_, ok := reservedProjectSlugs[pb.Slug]

	return ok
}
