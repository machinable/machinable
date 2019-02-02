// A stored Session is really a reference to a refresh token
package mongo

import (
	"context"
	"time"

	"github.com/anothrnick/machinable/dsi/models"
	"github.com/mongodb/mongo-go-driver/bson"
	"github.com/mongodb/mongo-go-driver/bson/objectid"
)

// CreateAppSession create new session for an application user
func (d *Datastore) CreateAppSession(session *models.Session) error {
	session.ID = objectid.New()

	col := d.db.Sessions()
	_, err := col.InsertOne(
		context.Background(),
		session,
	)

	return err
}

// UpdateAppSessionLastAccessed update session last accessed
func (d *Datastore) UpdateAppSessionLastAccessed(sessionID string, lastAccessed time.Time) error {
	col := d.db.Sessions()

	sessionObjectID, err := objectid.FromHex(sessionID)

	if err != nil {
		return err
	}

	// update session `last_accessed` time
	_, err = col.UpdateOne(
		context.Background(),
		bson.NewDocument(
			bson.EC.ObjectID("_id", sessionObjectID),
		),
		bson.NewDocument(
			bson.EC.SubDocumentFromElements("$set",
				bson.EC.Time("last_accessed", lastAccessed),
			),
		),
	)

	return err
}

// GetAppSession retrieve a single application session by ID
func (d *Datastore) GetAppSession(sessionID string) (*models.Session, error) {
	session := &models.Session{}
	col := d.db.Sessions()

	sessionObjectID, err := objectid.FromHex(sessionID)

	if err != nil {
		return nil, err
	}

	err = col.FindOne(
		nil,
		bson.NewDocument(
			bson.EC.ObjectID("_id", sessionObjectID),
		),
		nil,
	).Decode(session)

	return session, err
}

// DeleteAppSession permanently remove the session by ID
func (d *Datastore) DeleteAppSession(sessionID string) error {
	col := d.db.Sessions()

	sessionObjectID, err := objectid.FromHex(sessionID)

	if err != nil {
		return err
	}

	_, err = col.DeleteOne(
		context.Background(),
		bson.NewDocument(
			bson.EC.ObjectID("_id", sessionObjectID),
		),
	)

	return err
}
