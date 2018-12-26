package interfaces

import "bitbucket.org/nsjostrom/machinable/dsi/models"

// SessionsDatastore exposes functions to manage application user sessions
type SessionsDatastore interface {
	CreateAppSession(session *models.Session) error
	GetAppSession(sessionID string) (*models.Session, error)
	ListAppSessions() ([]*models.Session, error)
	DeleteAppSession(sessionID string) error
}
