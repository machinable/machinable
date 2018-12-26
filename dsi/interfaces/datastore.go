package interfaces

// Datastore exposes the necessary functions to interact with the Machinable datastore.
// Functions are grouped logically based on their purpose and the collections they interact with.
// implemented connectors: MongoDB
// potential connectors: InfluxDB, Postgres JSON, Redis, CouchDB, etc.
type Datastore interface {
	// Project resources/definitions
	ResourcesDatastore
	// Project collections
	CollectionsDatastore
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

}
