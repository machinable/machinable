package mongo

import "github.com/anothrnick/machinable/dsi/mongo/database"

// New returns a new MongoDatastore struct. `database.Database` is required.
func New(db *database.Database) *Datastore {
	return &Datastore{
		db: db,
	}
}

// Datastore is the mongoDB implementation of the Datastore interface
type Datastore struct {
	db *database.Database
}
