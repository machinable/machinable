package mongo

import (
	"context"
	"errors"
	"time"

	"bitbucket.org/nsjostrom/machinable/dsi/models"
	"github.com/mongodb/mongo-go-driver/bson"
	"github.com/mongodb/mongo-go-driver/bson/objectid"
)

// UpdateProjectAuthn updates the project authentication policy
func (d *Datastore) UpdateProjectAuthn(slug, userID string, authn bool) (*models.Project, error) {
	// Create object ID from resource ID string
	userObjectID, err := objectid.FromHex(userID)
	if err != nil {
		return nil, err
	}

	// connect to project collection
	col := d.db.Projects()

	// only updated `authn` value
	_, err = col.UpdateOne(
		context.Background(),
		bson.NewDocument(
			bson.EC.String("slug", slug),
			bson.EC.ObjectID("user_id", userObjectID),
		),
		bson.NewDocument(
			bson.EC.SubDocumentFromElements("$set",
				bson.EC.Boolean("authn", authn),
			),
		),
	)

	if err != nil {
		return nil, err
	}

	project, err := d.GetProjectBySlug(slug)

	return project, err
}

// CreateProject creates a new project for a user
func (d *Datastore) CreateProject(userID, slug, name, description, icon string, authn bool) (*models.Project, error) {
	// create ObjectID from UserID string
	userObjectID, err := objectid.FromHex(userID)
	if err != nil {
		return nil, err
	}

	// connect to project collection
	col := d.db.Projects()

	// init project struc
	project := &models.Project{
		UserID:      userObjectID,
		Slug:        slug,
		Name:        name,
		Description: description,
		Icon:        icon,
		Created:     time.Now(),
		Authn:       authn,
	}

	// save user project
	_, err = col.InsertOne(
		context.Background(),
		project,
	)

	return project, err
}

// GetProjectBySlug retrieves a project by slug
func (d *Datastore) GetProjectBySlug(slug string) (*models.Project, error) {
	// connect to project collection
	col := d.db.Projects()

	project := &models.Project{}
	// look up the user
	err := col.FindOne(
		nil,
		bson.NewDocument(
			bson.EC.String("slug", slug),
		),
		nil,
	).Decode(project)

	return project, err
}

// ListUserProjects retrieves all projects for a user
func (d *Datastore) ListUserProjects(userID string) ([]*models.Project, error) {
	// create ObjectID from UserID string
	userObjectID, err := objectid.FromHex(userID)
	if err != nil {
		return nil, err
	}

	// connect to project collection
	col := d.db.Projects()

	// look up projects
	cursor, err := col.Find(
		nil,
		bson.NewDocument(
			bson.EC.ObjectID("user_id", userObjectID),
		),
	)
	if err != nil {
		return nil, err
	}

	projects := make([]*models.Project, 0)
	for cursor.Next(context.Background()) {
		prj := &models.Project{}
		err := cursor.Decode(prj)
		if err != nil {
			return nil, err
		}
		projects = append(projects, prj)
	}

	return projects, err
}

// DeleteProject permanently removes a project based on project slug
func (d *Datastore) DeleteProject(slug string) error {
	return errors.New("not implemented")
}
