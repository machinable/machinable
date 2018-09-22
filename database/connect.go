package database

import (
	"context"
	"fmt"
	"os"

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

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

// Connect returns a *mongo.Database connection
func Connect() *mongo.Database {
	host := getEnv("MONGO_HOST", "localhost")
	port := getEnv("MONGO_PORT", "27017")
	client, err := mongo.Connect(context.Background(), fmt.Sprintf("mongodb://%s:%s", host, port), nil)
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
