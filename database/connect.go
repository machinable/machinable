package database

import (
	"context"
	"fmt"

	"github.com/mongodb/mongo-go-driver/mongo"
)

const (
	// DatabaseName is the name of the application database
	DatabaseName = "flowdb"
	// ResourceDefinitions is the collection for storing resource definitions
	ResourceDefinitions = "definitions"
	// ResourceFormat is the string format for a resource, this should include an account specifier as well
	ResourceFormat = "resource.%s"
)

// Connect returns a *mongo.Database connection
func Connect() *mongo.Database {
	client, err := mongo.Connect(context.Background(), "mongodb://localhost:27017", nil)
	if err != nil {
		fmt.Println(err)
		panic("failed to connect database")
	}

	return client.Database(DatabaseName)
}

// Collection returns a *mongo.Collection connection
func Collection(col string) *mongo.Collection {
	return Connect().Collection(col)
}
