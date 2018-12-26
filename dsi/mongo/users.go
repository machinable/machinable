package mongo

import (
	"context"

	"bitbucket.org/nsjostrom/machinable/dsi/models"
	"github.com/mongodb/mongo-go-driver/bson"
	"github.com/mongodb/mongo-go-driver/bson/objectid"
)

// GetAppUserByUsername attempts to find a user by username, if the user does not exist the user will be nil, error will be !nil
func (d *Datastore) GetAppUserByUsername(userName string) (*models.User, error) {
	col := d.db.Users()

	user := &models.User{}
	err := col.FindOne(
		nil,
		bson.NewDocument(
			bson.EC.String("username", userName),
		),
		nil,
	).Decode(user)

	if err != nil {
		return nil, err
	}

	return user, nil
}

// GetAppUserByID attempts to find a user by ID, if the user does not exist the user will be nil, error will be !nil
func (d *Datastore) GetAppUserByID(id string) (*models.User, error) {
	userObjectID, err := objectid.FromHex(id)
	if err != nil {
		return nil, err
	}

	col := d.db.Users()

	user := &models.User{}
	err = col.FindOne(
		nil,
		bson.NewDocument(
			bson.EC.ObjectID("_id", userObjectID),
		),
		nil,
	).Decode(user)

	if err != nil {
		return nil, err
	}

	return user, nil
}

// CreateAppUser saves a new application user
func (d *Datastore) CreateAppUser(user *models.User) error {
	user.ID = objectid.New()
	col := d.db.Users()

	// save the user
	_, err := col.InsertOne(
		context.Background(),
		user,
	)

	return err
}
