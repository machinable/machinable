package database

import (
	"context"
	"fmt"
	"os"

	"github.com/mongodb/mongo-go-driver/mongo"
)

const (
	databaseName        = "machinable"
	resourceDefinitions = "project.%s.definitions"
	resourceFormat      = "project.%s.resource.%s"
	collections         = "project.%s.collections"
	collectionFormat    = "project.%s.collections.%s"
	userCollection      = "project.%s.users"
	tokenCollection     = "project.%s.keys"
	sessionCollection   = "project.%s.sessions"
	logCollection       = "project.%s.logs"

	usersCollection    = "users"
	sessionsCollection = "sessions"
	projectsCollection = "projects"
)

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

func createConnection() *mongo.Database {
	host := getEnv("MONGO_HOST", "localhost")
	port := getEnv("MONGO_PORT", "27017")
	client, err := mongo.Connect(context.Background(), fmt.Sprintf("mongodb://%s:%s", host, port), nil)
	if err != nil {
		fmt.Println(err)
		panic("failed to connect database")
	}

	return client.Database(databaseName)
}

// Connect returns a *mongo.Database connection
func Connect() *Database {
	mongoDatabase := createConnection()

	return &Database{db: mongoDatabase}
}

// Database wraps the mongo database connection and provides helper functions for collections
type Database struct {
	db *mongo.Database
}

// Projects returns the app projects collection
func (d *Database) Projects() *mongo.Collection {
	return d.db.Collection(projectsCollection)
}

// Users returns the app users collection
func (d *Database) Users() *mongo.Collection {
	return d.db.Collection(usersCollection)
}

// Sessions returns the app sessions collection
func (d *Database) Sessions() *mongo.Collection {
	return d.db.Collection(sessionsCollection)
}

// LogDocs the collection of project logs
func (d *Database) LogDocs(projectSlug string) *mongo.Collection {
	return d.db.Collection(fmt.Sprintf(logCollection, projectSlug))
}

// SessionDocs the collection of project sessions
func (d *Database) SessionDocs(projectSlug string) *mongo.Collection {
	return d.db.Collection(fmt.Sprintf(sessionCollection, projectSlug))
}

// KeyDocs returns the collection that stores project api Keys
func (d *Database) KeyDocs(projectSlug string) *mongo.Collection {
	return d.db.Collection(fmt.Sprintf(tokenCollection, projectSlug))
}

// UserDocs returns the collection that stores project users
func (d *Database) UserDocs(projectSlug string) *mongo.Collection {
	return d.db.Collection(fmt.Sprintf(userCollection, projectSlug))
}

// ResourceDefinitions returns the collection that stores resource definitons for a project
func (d *Database) ResourceDefinitions(projectSlug string) *mongo.Collection {
	return d.db.Collection(fmt.Sprintf(resourceDefinitions, projectSlug))
}

// ResourceDocs returns the collection that stores resources (documents) for a project for a resource (path name)
func (d *Database) ResourceDocs(projectSlug, resourcePath string) *mongo.Collection {
	return d.db.Collection(fmt.Sprintf(resourceFormat, projectSlug, resourcePath))
}

// CollectionNames returns the formatted string of the collection name of the collection that stores the list of project collections
func (d *Database) CollectionNames(projectSlug string) *mongo.Collection {
	return d.db.Collection(fmt.Sprintf(collections, projectSlug))
}

// CollectionDocs returns the collection that stores the list of documents for a project collection
func (d *Database) CollectionDocs(projectSlug, collectionName string) *mongo.Collection {
	return d.db.Collection(fmt.Sprintf(collectionFormat, projectSlug, collectionName))
}
