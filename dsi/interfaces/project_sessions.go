package interfaces

import (
	"time"

	"github.com/anothrnick/machinable/dsi/models"
)

// ProjectSessionsDatastore exposes functions to manage project user sessions
type ProjectSessionsDatastore interface {
	CreateSession(projectID string, session *models.Session) error
	UpdateProjectSessionLastAccessed(projectID, sessionID string, lastAccessed time.Time) error
	GetSession(projectID, sessionID string) (*models.Session, error)
	ListSessions(projectID string) ([]*models.Session, error)
	DeleteSession(projectID, sessionID string) error
	DropProjectSessions(projectID string) error
}
