package interfaces

import "github.com/anothrNick/machinable/dsi/models"

// Datastore exposes the necessary functions to interact with the Machinable datastore.
// Functions are grouped logically based on their purpose and the collections they interact with.
// implemented connectors: MongoDB
// potential connectors: InfluxDB, Postgres JSON, Redis, CouchDB, etc.
type Datastore interface {
	// Project resources/definitions
	ResourcesDatastore
	// JSON Key/val
	ProjectJSONDatastore
	// Project users
	ProjectUsersDatastore
	// Project apikeys
	ProjectAPIKeysDatastore
	// Project logs
	ProjectLogsDatastore
	// Project sessions
	ProjectSessionsDatastore
	// Projects
	ProjectsDatastore
	// Users
	UsersDatastore
	// Sessions
	SessionsDatastore
	// Tiers
	TiersDatastore
	// Errors
	TranslateError(err error) *models.TranslatedError
}
