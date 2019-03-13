package interfaces

import (
	"time"

	"github.com/anothrnick/machinable/dsi/models"
)

// ProjectSessionsDatastore exposes functions to manage project user sessions
type ProjectSessionsDatastore interface {
	CreateSession(project string, session *models.Session) error
	UpdateProjectSessionLastAccessed(project, sessionID string, lastAccessed time.Time) error
	GetSession(project, sessionID string) (*models.Session, error)
	ListSessions(project string) ([]*models.Session, error)
	DeleteSession(project, sessionID string) error
	DropProjectSessions(project string) error
}
