package mongo

import (
	"context"

	"github.com/anothrnick/machinable/dsi/models"
	"github.com/mongodb/mongo-go-driver/bson"
	"github.com/mongodb/mongo-go-driver/bson/objectid"
)

// GetUserByUsername retrieves a project user by the user's username
func (d *Datastore) GetUserByUsername(project, userName string) (*models.ProjectUser, error) {
	// get the users collection
	collection := d.db.UserDocs(project)

	// look up the user
	documentResult := collection.FindOne(
		nil,
		bson.NewDocument(
			bson.EC.String("username", userName),
		),
		nil,
	)

	user := &models.ProjectUser{}
	// decode user document
	err := documentResult.Decode(user)

	return user, err
}

// CreateUser creates a new project user for the project
func (d *Datastore) CreateUser(project string, user *models.ProjectUser) error {
	// get the users collection
	col := d.db.UserDocs(project)
	// save the user
	_, err := col.InsertOne(
		context.Background(),
		user,
	)

	return err
}

// ListUsers returns all project users for a project
func (d *Datastore) ListUsers(project string) ([]*models.ProjectUser, error) {
	users := make([]*models.ProjectUser, 0)

	col := d.db.UserDocs(project)

	cursor, err := col.Find(
		context.Background(),
		bson.NewDocument(),
	)

	if err != nil {
		return users, err
	}

	for cursor.Next(context.Background()) {
		var user models.ProjectUser
		err := cursor.Decode(&user)
		if err != nil {
			return users, err
		}
		users = append(users, &user)
	}

	return users, nil
}

// DeleteUser deletes a project user for a project based on userID
func (d *Datastore) DeleteUser(project, userID string) error {
	// Create object ID from resource ID string
	objectID, err := objectid.FromHex(userID)
	if err != nil {
		return err
	}

	sessCollection := d.db.SessionDocs(project)
	// Delete the sessions
	_, err = sessCollection.DeleteMany(
		context.Background(),
		bson.NewDocument(
			bson.EC.String("user_id", userID),
		),
	)
	if err != nil {
		return err
	}

	userCollection := d.db.UserDocs(project)
	// Delete the user
	_, err = userCollection.DeleteOne(
		context.Background(),
		bson.NewDocument(
			bson.EC.ObjectID("_id", objectID),
		),
	)

	return err
}
