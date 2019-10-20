package postgres

import (
	"fmt"
	"strings"

	"github.com/anothrnick/machinable/dsi/models"
)

const tableProjectJSON = "project_json"

// GetRootKey retrieves a single root key by the key name
func (d *Database) GetRootKey(projectID, rootKey string) (*models.RootKey, error) {
	newKey := models.RootKey{}
	err := d.db.QueryRow(
		fmt.Sprintf(
			"SELECT id, project_id, root_key, \"create\", \"read\", \"update\", \"delete\" FROM %s WHERE project_id=$1 and root_key=$2",
			tableProjectJSON,
		),
		projectID,
		rootKey,
	).Scan(
		&newKey.ID,
		&newKey.ProjectID,
		&newKey.Key,
		&newKey.Create,
		&newKey.Read,
		&newKey.Update,
		&newKey.Delete,
	)

	return &newKey, err
}

// ListRootKeys lists all root keys with associated metadata, does not include jsonb
func (d *Database) ListRootKeys(projectID string) ([]*models.RootKey, error) {
	rows, err := d.db.Query(
		fmt.Sprintf(
			"SELECT id, project_id, root_key, \"create\", \"read\", \"update\", \"delete\" FROM %s WHERE project_id=$1",
			tableProjectJSON,
		),
		projectID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	rootKeys := make([]*models.RootKey, 0)
	for rows.Next() {
		rootKey := models.RootKey{}
		err = rows.Scan(
			&rootKey.ID,
			&rootKey.ProjectID,
			&rootKey.Key,
			&rootKey.Create,
			&rootKey.Read,
			&rootKey.Update,
			&rootKey.Delete,
		)
		if err != nil {
			return nil, err
		}

		rootKeys = append(rootKeys, &rootKey)
	}

	return rootKeys, rows.Err()
}

// CreateRootKey creates a new rootkey with a JSON tree
func (d *Database) CreateRootKey(projectID, rootKey string, data []byte) error {
	_, err := d.db.Exec(
		fmt.Sprintf(
			"INSERT INTO %s (project_id, root_key, data) VALUES ($1, $2, $3)",
			tableProjectJSON,
		),
		projectID,
		rootKey,
		data,
	)
	return err
}

//DeleteRootKey permanently deletes an entire rootkey's tree
func (d *Database) DeleteRootKey(projectID, rootKey string) error {
	_, err := d.db.Exec(
		fmt.Sprintf(
			"DELETE FROM %s where project_id=$1 and root_key=$2",
			tableProjectJSON,
		),
		projectID,
		rootKey,
	)
	return err
}

// GetJSONKey retrieves the object at the key path
func (d *Database) GetJSONKey(projectID, rootKey string, keys ...string) ([]byte, error) {
	byt := []byte{}

	// escape?
	keysFormat := strings.Join(keys, ",")
	err := d.db.QueryRow(
		fmt.Sprintf(
			"SELECT data#>'{%s}' as data FROM %s WHERE project_id=$1 and root_key=$2",
			keysFormat,
			tableProjectJSON,
		),
		projectID,
		rootKey,
	).Scan(&byt)

	return byt, err
}

// CreateJSONKey saves the data at the provided key path. Fails if the key already exists.
func (d *Database) CreateJSONKey(projectID, rootKey string, data []byte, keys ...string) error {
	// escape?
	keysFormat := strings.Join(keys, ",")
	_, err := d.db.Exec(
		fmt.Sprintf(
			"UPDATE %s set data=jsonb_insert(data, '{%s}', $1) WHERE project_id=$2 and root_key=$3",
			tableProjectJSON,
			keysFormat,
		),
		data,
		projectID,
		rootKey,
	)
	return err
}

// UpdateJSONKey updates the data at the key path. Creates a new key if it does not already exist.
func (d *Database) UpdateJSONKey(projectID, rootKey string, data []byte, keys ...string) error {
	// escape?
	keysFormat := strings.Join(keys, ",")
	_, err := d.db.Exec(
		fmt.Sprintf(
			"UPDATE %s set data=jsonb_set(data, '{%s}', $1) WHERE project_id=$2 and root_key=$3",
			tableProjectJSON,
			keysFormat,
		),
		data,
		projectID,
		rootKey,
	)
	return err
}

// DeleteJSONKey permanently removes the data at the key path.
func (d *Database) DeleteJSONKey(projectID, rootKey string, keys ...string) error {
	// escape?
	keysFormat := strings.Join(keys, ",")
	_, err := d.db.Exec(
		fmt.Sprintf(
			"UPDATE %s SET data=data #- '{%s}' WHERE project_id=$1 and root_key=$2",
			tableProjectJSON,
			keysFormat,
		),
		projectID,
		rootKey,
	)
	return err
}
