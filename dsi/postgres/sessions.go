package postgres

import (
	"time"

	"github.com/anothrnick/machinable/dsi/models"
)

const tableAppSessions = "app_sessions"

// CreateAppSession create new session for an application user
func (d *Database) CreateAppSession(session *models.Session) error {
	return nil
}

// UpdateAppSessionLastAccessed update session last accessed
func (d *Database) UpdateAppSessionLastAccessed(sessionID string, lastAccessed time.Time) error {
	return nil
}

// ListUserSessions lists all sessions for a user
func (d *Database) ListUserSessions(userID string) ([]*models.Session, error) {
	return nil, nil
}

// GetAppSession retrieve a single application session by ID
func (d *Database) GetAppSession(sessionID string) (*models.Session, error) {
	return nil, nil
}

// DeleteAppSession permanently remove the session by ID
func (d *Database) DeleteAppSession(sessionID string) error {
	return nil
}
