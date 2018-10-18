package database

import (
	"context"
	"fmt"
	"os"

	"github.com/mongodb/mongo-go-driver/mongo"
)

const (
	databaseName = "machinable"
	// Users of the app
	Users = "users"
	// Sessions are active user sessions which represent refresh tokens for the users' JWT
	Sessions = "sessions"
	// Projects are user projects
	Projects = "projects"
)

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
