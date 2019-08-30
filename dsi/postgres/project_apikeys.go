package postgres

import (
	"fmt"
	"time"

	"github.com/anothrnick/machinable/dsi/models"
)

const tableProjectAPIKeys = "project_apikeys"

// GetAPIKeyByKey retrieves a single api key by key hash
func (d *Database) GetAPIKeyByKey(projectID, hash string) (*models.ProjectAPIKey, error) {
	key := models.ProjectAPIKey{}
	err := d.db.QueryRow(
		fmt.Sprintf(
			"SELECT id, project_id, key_hash, read, write, role, created FROM %s WHERE project_id=$1 and key_hash=$2",
			tableProjectAPIKeys,
		),
		projectID,
		hash,
	).Scan(
		&key.ID,
		&key.ProjectID,
		&key.KeyHash,
		&key.Read,
		&key.Write,
		&key.Role,
		&key.Created,
	)

	return &key, err
}

// CreateAPIKey creates a new api key for the project
func (d *Database) CreateAPIKey(projectID, hash, description string, read, write bool, role string) (*models.ProjectAPIKey, error) {
	_, err := d.db.Exec(
		fmt.Sprintf(
			"INSERT INTO %s (project_id, key_hash, description, read, write, role, created) VALUES ($1, $2, $3, $4, $5, $6, $7)",
			tableProjectAPIKeys,
		),
		projectID,
		hash,
		description,
		read,
		write,
		role,
		time.Now(),
	)

	return nil, err
}

// UpdateAPIKey updates the role and access of an API key
func (d *Database) UpdateAPIKey(projectID, keyID string, read, write bool, role string) error {
	_, err := d.db.Exec(
		fmt.Sprintf(
			"UPDATE %s SET read=$1, write=$2, role=$3 WHERE id=$4 and project_id=$5",
			tableProjectAPIKeys,
		),
		read,
		write,
		role,
		keyID,
		projectID,
	)

	return err
}

// ListAPIKeys retrieves all api keys for a project
func (d *Database) ListAPIKeys(projectID string) ([]*models.ProjectAPIKey, error) {
	rows, err := d.db.Query(
		fmt.Sprintf(
			"SELECT id, project_id, key_hash, read, write, role, created FROM %s WHERE project_id=$1",
			tableProjectAPIKeys,
		),
		projectID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	keys := make([]*models.ProjectAPIKey, 0)
	for rows.Next() {
		key := models.ProjectAPIKey{}
		err = rows.Scan(
			&key.ID,
			&key.ProjectID,
			&key.KeyHash,
			&key.Read,
			&key.Write,
			&key.Role,
			&key.Created,
		)
		if err != nil {
			return nil, err
		}

		keys = append(keys, &key)
	}

	return keys, rows.Err()
}

// DeleteAPIKey removes a project api key permanently
func (d *Database) DeleteAPIKey(projectID, keyID string) error {
	_, err := d.db.Exec(
		fmt.Sprintf(
			"DELETE FROM %s WHERE id=$1 and project_id=$2",
			tableProjectAPIKeys,
		),
		keyID,
		projectID,
	)
	return err
}

// DropProjectKeys drops the key collection for this project
func (d *Database) DropProjectKeys(projectID string) error {
	_, err := d.db.Exec(
		fmt.Sprintf(
			"DELETE FROM %s WHERE project_id=$1",
			tableProjectAPIKeys,
		),
		projectID,
	)
	return err
}
