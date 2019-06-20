package interfaces

import (
	"time"

	"github.com/anothrnick/machinable/dsi/models"
)

// SessionsDatastore exposes functions to manage application user sessions
type SessionsDatastore interface {
	CreateAppSession(session *models.Session) error
	UpdateAppSessionLastAccessed(sessionID string, lastAccessed time.Time) error
	ListUserSessions(userID string) ([]*models.Session, error)
	GetAppSession(sessionID string) (*models.Session, error)
	DeleteAppSession(sessionID string) error
}
