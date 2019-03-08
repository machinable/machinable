package mongo

import (
	"context"

	"github.com/anothrnick/machinable/dsi/models"
	"github.com/mongodb/mongo-go-driver/bson"
	"github.com/mongodb/mongo-go-driver/bson/objectid"
)

// CreateSession creates a new session for a project user
func (d *Datastore) CreateSession(project string, session *models.Session) error {
	collection := d.db.SessionDocs(project)
	// save the session
	_, err := collection.InsertOne(
		context.Background(),
		session,
	)

	return err
}

// GetSession retrieves a single project session by ID
func (d *Datastore) GetSession(project, sessionID string) (*models.Session, error) {
	return nil, nil
}

// ListSessions lists all sessions for a project
func (d *Datastore) ListSessions(project string) ([]*models.Session, error) {
	sessions := make([]*models.Session, 0)
	collection := d.db.SessionDocs(project)

	cursor, err := collection.Find(
		context.Background(),
		bson.NewDocument(),
	)

	if err != nil {
		return nil, err
	}

	for cursor.Next(context.Background()) {
		var session models.Session
		err := cursor.Decode(&session)
		if err != nil {
			return nil, err
		}
		sessions = append(sessions, &session)
	}

	return sessions, nil
}

// DeleteSession removes a project user's session by project and ID
func (d *Datastore) DeleteSession(project, sessionID string) error {
	collection := d.db.SessionDocs(project)
	// Get the object id
	objectID, err := objectid.FromHex(sessionID)
	if err != nil {
		return err
	}

	// Delete the session
	_, err = collection.DeleteOne(
		context.Background(),
		bson.NewDocument(
			bson.EC.ObjectID("_id", objectID),
		),
	)

	return err
}

// DropProjectSessions drops the collection of this project's user sessions
func (d *Datastore) DropProjectSessions(project string) error {
	// drop collection storage
	collection := d.db.SessionDocs(project)

	err := collection.Drop(nil, nil)

	return err
}
