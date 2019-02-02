package interfaces

import "github.com/anothrnick/machinable/dsi/models"

// ProjectSessionsDatastore exposes functions to manage project user sessions
type ProjectSessionsDatastore interface {
	CreateSession(project string, session *models.Session) error
	GetSession(project, sessionID string) (*models.Session, error)
	ListSessions(project string) ([]*models.Session, error)
	DeleteSession(project, sessionID string) error
}
