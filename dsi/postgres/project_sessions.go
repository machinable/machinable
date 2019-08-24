package postgres

import (
	"time"

	"github.com/anothrnick/machinable/dsi/models"
)

const tableProjectSessions = "project_sessions"

// CreateSession creates a new session for a project user
func (d *Database) CreateSession(project string, session *models.Session) error {
	return nil
}

// UpdateProjectSessionLastAccessed update session last accessed
func (d *Database) UpdateProjectSessionLastAccessed(project, sessionID string, lastAccessed time.Time) error {
	return nil
}

// GetSession retrieves a single project session by ID
func (d *Database) GetSession(project, sessionID string) (*models.Session, error) {
	return nil, nil
}

// ListSessions lists all sessions for a project
func (d *Database) ListSessions(project string) ([]*models.Session, error) {
	return nil, nil
}

// DeleteSession removes a project user's session by project and ID
func (d *Database) DeleteSession(project, sessionID string) error {
	return nil
}

// DropProjectSessions drops the collection of this project's user sessions
func (d *Database) DropProjectSessions(project string) error {
	return nil
}
