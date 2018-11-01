package database

import (
	"context"
	"fmt"
	"os"

	"github.com/mongodb/mongo-go-driver/mongo"
)

const (
	databaseName        = "machinable"
	resourceDefinitions = "%s.definitions"
	resourceFormat      = "%s.resource.%s"
	collections         = "%s.collections"
	collectionFormat    = "%s.collections.%s"
	userCollection      = "%s.users"
	tokenCollection     = "%s.keys"
	sessionCollection   = "%s.sessions"
	logCollection       = "%s.logs"
)

// LogDocs just that
func LogDocs(projectSlug string) string {
	return fmt.Sprintf(logCollection, projectSlug)
}

// SessionDocs returns the formatted string of the collection name for the collection of sessions
func SessionDocs(projectSlug string) string {
	return fmt.Sprintf(sessionCollection, projectSlug)
}

// KeyDocs returns the formatted string of the collection name of the collection that stores project api Keys
func KeyDocs(projectSlug string) string {
	return fmt.Sprintf(tokenCollection, projectSlug)
}

// UserDocs returns the formatted string of the collection name of the collection that stores project users
func UserDocs(projectSlug string) string {
	return fmt.Sprintf(userCollection, projectSlug)
}

// ResourceDefinitions returns the formatted string of the collection name of the collection that stores resource definitons for a project
func ResourceDefinitions(projectSlug string) string {
	return fmt.Sprintf(resourceDefinitions, projectSlug)
}

// ResourceDocs returns the formatted string of the collection name of the collection that stores resources (documents) for a project for a resource (path name)
func ResourceDocs(projectSlug, resourcePath string) string {
	return fmt.Sprintf(resourceFormat, projectSlug, resourcePath)
}

// CollectionNames returns the formatted string of the collection name of the collection that stores the list of project collections
func CollectionNames(projectSlug string) string {
	return fmt.Sprintf(collections, projectSlug)
}

// CollectionDocs returns the formatted string of the collection name of the collection that stores the list of documents for a project collection
func CollectionDocs(projectSlug, collectionName string) string {
	return fmt.Sprintf(collectionFormat, projectSlug, collectionName)
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

// DB is the connectin to the mongodb, TODO: move this to a managed structure
var DB *mongo.Database

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
func Connect() *mongo.Database {
	if DB == nil {
		DB = createConnection()
	}

	return DB
}

// Collection returns a *mongo.Collection connection
func Collection(col string) *mongo.Collection {
	return Connect().Collection(col)
}
