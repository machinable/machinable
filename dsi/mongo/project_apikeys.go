package mongo

import (
	"context"
	"time"

	"github.com/anothrnick/machinable/dsi/models"
	"github.com/mongodb/mongo-go-driver/bson"
	"github.com/mongodb/mongo-go-driver/bson/objectid"
)

// GetAPIKeyByKey retrieves a single api key by key hash
func (d *Datastore) GetAPIKeyByKey(project, hash string) (*models.ProjectAPIKey, error) {
	// get the keys collection
	col := d.db.KeyDocs(project)
	key := &models.ProjectAPIKey{}

	// Find api key
	err := col.FindOne(
		nil,
		bson.NewDocument(
			bson.EC.String("key_hash", hash),
		),
		nil,
	).Decode(key)

	return key, err
}

// CreateAPIKey creates a new api key for the project
func (d *Datastore) CreateAPIKey(project, hash, description string, read, write bool, role string) (*models.ProjectAPIKey, error) {
	key := &models.ProjectAPIKey{
		ID:          objectid.New(), // I don't like this
		Created:     time.Now(),
		KeyHash:     hash,
		Description: description,
		Read:        read,
		Write:       write,
		Role:        role,
	}

	// get the keys collection
	col := d.db.KeyDocs(project)
	// save key
	_, err := col.InsertOne(
		context.Background(),
		key,
	)

	return key, err
}

// ListAPIKeys retrieves all api keys for a project
func (d *Datastore) ListAPIKeys(project string) ([]*models.ProjectAPIKey, error) {
	// get the keys collection
	col := d.db.KeyDocs(project)

	cursor, err := col.Find(
		context.Background(),
		bson.NewDocument(),
	)

	if err != nil {
		return nil, err
	}

	keys := make([]*models.ProjectAPIKey, 0)
	for cursor.Next(context.Background()) {
		var key models.ProjectAPIKey
		err := cursor.Decode(&key)
		if err != nil {
			return nil, err
		}
		keys = append(keys, &key)
	}

	return keys, err
}

// DeleteAPIKey removes a project api key permanently
func (d *Datastore) DeleteAPIKey(project, keyID string) error {
	// Create object ID from resource ID string
	objectID, err := objectid.FromHex(keyID)
	if err != nil {
		return err
	}

	// get the keys collection
	col := d.db.KeyDocs(project)

	// Delete the object
	_, err = col.DeleteOne(
		context.Background(),
		bson.NewDocument(
			bson.EC.ObjectID("_id", objectID),
		),
	)

	return err
}

// DropProjectKeys drops the key collection for this project
func (d *Datastore) DropProjectKeys(project string) error {
	collection := d.db.KeyDocs(project)

	err := collection.Drop(nil, nil)

	return err
}
