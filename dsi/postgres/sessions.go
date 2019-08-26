package postgres

import (
	"fmt"
	"time"

	"github.com/anothrnick/machinable/dsi/models"
)

const tableAppSessions = "app_sessions"

// CreateAppSession create new session for an application user
func (d *Database) CreateAppSession(session *models.Session) error {
	_, err := d.db.Exec(
		fmt.Sprintf(
			"INSERT INTO %s (user_id, location, mobile, ip, last_accessed, browser, os) VALUES ($1, $2, $3, $4, $5, $6, $7)",
			tableAppSessions,
		),
		session.UserID,
		session.Location,
		session.Mobile,
		session.IP,
		session.LastAccessed,
		session.Browser,
		session.OS,
	)

	return err
}

// UpdateAppSessionLastAccessed update session last accessed
func (d *Database) UpdateAppSessionLastAccessed(sessionID string, lastAccessed time.Time) error {
	_, err := d.db.Exec(
		fmt.Sprintf(
			"UPDATE %s SET last_accessed=$1 WHERE id=$2",
			tableAppSessions,
		),
		lastAccessed,
		sessionID,
	)

	return err
}

// ListUserSessions lists all sessions for a user
func (d *Database) ListUserSessions(userID string) ([]*models.Session, error) {
	rows, err := d.db.Query(
		fmt.Sprintf(
			"SELECT id, user_id, location, mobile, ip, last_accessed, browser, os FROM %s WHERE user_id=$1",
			tableAppSessions,
		),
		userID,
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

// GetAppSession retrieve a single application session by ID
func (d *Database) GetAppSession(sessionID string) (*models.Session, error) {
	session := models.Session{}

	err := d.db.QueryRow(
		fmt.Sprintf(
			"SELECT id, user_id, location, mobile, ip, last_accessed, browser, os FROM %s WHERE id=$1",
			tableAppSessions,
		),
		sessionID,
	).Scan(
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

// DeleteAppSession permanently remove the session by ID
func (d *Database) DeleteAppSession(sessionID string) error {
	_, err := d.db.Exec(
		fmt.Sprintf(
			"DELETE FROM %s WHERE id=$1",
			tableAppSessions,
		),
		sessionID,
	)
	return err
}
