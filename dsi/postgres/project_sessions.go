package postgres

import (
	"fmt"
	"time"

	"github.com/machinable/machinable/dsi/models"
)

const tableProjectSessions = "project_sessions"

// CreateSession creates a new session for a project user
func (d *Database) CreateSession(projectID string, session *models.Session) error {
	err := d.db.QueryRow(
		fmt.Sprintf(
			"INSERT INTO %s (project_id, user_id, location, mobile, ip, last_accessed, browser, os) VALUES ($1, $2, $3, $4, $5, $6, $7, $8) RETURNING id",
			tableProjectSessions,
		),
		projectID,
		session.UserID,
		session.Location,
		session.Mobile,
		session.IP,
		session.LastAccessed,
		session.Browser,
		session.OS,
	).Scan(&session.ID)

	return err
}

// UpdateProjectSessionLastAccessed update session last accessed
func (d *Database) UpdateProjectSessionLastAccessed(projectID, sessionID string, lastAccessed time.Time) error {
	_, err := d.db.Exec(
		fmt.Sprintf(
			"UPDATE %s SET last_accessed=$1 WHERE id=$2 and project_id=$3",
			tableProjectSessions,
		),
		lastAccessed,
		sessionID,
		projectID,
	)

	return err
}

// GetSession retrieves a single project session by ID
func (d *Database) GetSession(projectID, sessionID string) (*models.Session, error) {
	session := models.Session{}

	err := d.db.QueryRow(
		fmt.Sprintf(
			"SELECT id, project_id, user_id, location, mobile, ip, last_accessed, browser, os FROM %s WHERE id=$1 and project_id=$2",
			tableProjectSessions,
		),
		sessionID,
		projectID,
	).Scan(
		&session.ID,
		&session.ProjectID,
		&session.UserID,
		&session.Location,
		&session.Mobile,
		&session.IP,
		&session.LastAccessed,
		&session.Browser,
		&session.OS,
	)
	if err != nil {
		return nil, err
	}

	return &session, err
}

// ListSessions lists all sessions for a project
func (d *Database) ListSessions(projectID string) ([]*models.Session, error) {
	rows, err := d.db.Query(
		fmt.Sprintf(
			"SELECT id, user_id, project_id, location, mobile, ip, last_accessed, browser, os FROM %s WHERE project_id=$1",
			tableProjectSessions,
		),
		projectID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	sessions := make([]*models.Session, 0)
	for rows.Next() {
		session := models.Session{}
		err = rows.Scan(
			&session.ID,
			&session.UserID,
			&session.ProjectID,
			&session.Location,
			&session.Mobile,
			&session.IP,
			&session.LastAccessed,
			&session.Browser,
			&session.OS,
		)
		if err != nil {
			return nil, err
		}

		sessions = append(sessions, &session)
	}

	return sessions, rows.Err()
}

// DeleteSession removes a project user's session by project and ID
func (d *Database) DeleteSession(projectID, sessionID string) error {
	_, err := d.db.Exec(
		fmt.Sprintf(
			"DELETE FROM %s WHERE id=$1 and project_id=$2",
			tableProjectSessions,
		),
		sessionID,
		projectID,
	)
	return err
}

// DropProjectSessions drops the collection of this project's user sessions
func (d *Database) DropProjectSessions(projectID string) error {
	_, err := d.db.Exec(
		fmt.Sprintf(
			"DELETE FROM %s WHERE project_id=$1",
			tableProjectSessions,
		),
		projectID,
	)
	return err
}
