package postgres

import (
	"fmt"
	"time"

	"github.com/anothrnick/machinable/dsi/models"
)

const tableProjectUsers = "project_users"

// GetUserByUsername retrieves a project user by the user's username
func (d *Database) GetUserByUsername(projectID, userName string) (*models.ProjectUser, error) {
	user := &models.ProjectUser{}

	err := d.db.QueryRow(
		fmt.Sprintf(
			"SELECT id, project_id, email, username, password_hash, read, write, role, created from %s WHERE username=$1 and project_id=$2",
			tableProjectUsers,
		),
		userName,
		projectID,
	).Scan(
		&user.ID,
		&user.ProjectID,
		&user.Email,
		&user.Username,
		&user.PasswordHash,
		&user.Read,
		&user.Write,
		&user.Role,
		&user.Created,
	)

	return user, err
}

// GetUserByID retrieves a project user by user _id
func (d *Database) GetUserByID(projectID, userID string) (*models.ProjectUser, error) {
	user := &models.ProjectUser{}

	err := d.db.QueryRow(
		fmt.Sprintf(
			"SELECT id, project_id, email, username, password_hash, read, write, role, created from %s WHERE id=$1 and project_id=$2",
			tableProjectUsers,
		),
		userID,
		projectID,
	).Scan(
		&user.ID,
		&user.Email,
		&user.Username,
		&user.PasswordHash,
		&user.Read,
		&user.Write,
		&user.Role,
		&user.Created,
	)

	return user, err
}

// CreateUser creates a new project user for the project
func (d *Database) CreateUser(projectID string, user *models.ProjectUser) error {
	_, err := d.db.Exec(
		fmt.Sprintf(
			"INSERT INTO %s (project_id, email, username, password_hash, read, write, role, created) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)",
			tableProjectUsers,
		),
		projectID,
		user.Email,
		user.Username,
		user.PasswordHash,
		user.Read,
		user.Write,
		user.Role,
		time.Now(),
	)

	return err
}

// UpdateUser updates the project user's access and role
func (d *Database) UpdateUser(projectID, userID string, user *models.ProjectUser) error {
	_, err := d.db.Exec(
		fmt.Sprintf(
			"UPDATE %s SET read=$1, write=$2, role=$3 WHERE id=$4 and project_id=$5",
			tableProjectUsers,
		),
		user.Read,
		user.Write,
		user.Role,
		userID,
		projectID,
	)

	return err
}

// ListUsers returns all project users for a project
func (d *Database) ListUsers(projectID string) ([]*models.ProjectUser, error) {
	rows, err := d.db.Query(
		fmt.Sprintf(
			"SELECT id, project_id, email, username, password_hash, read, write, role, created FROM %s WHERE project_id=$1",
			tableProjectUsers,
		),
		projectID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	users := make([]*models.ProjectUser, 0)
	for rows.Next() {
		user := models.ProjectUser{}
		err = rows.Scan(
			&user.ID,
			&user.ProjectID,
			&user.Email,
			&user.Username,
			&user.PasswordHash,
			&user.Read,
			&user.Write,
			&user.Role,
			&user.Created,
		)
		if err != nil {
			return nil, err
		}

		users = append(users, &user)
	}

	return users, rows.Err()
}

// DeleteUser deletes a project user for a project based on userID
func (d *Database) DeleteUser(projectID, userID string) error {
	_, err := d.db.Exec(
		fmt.Sprintf(
			"DELETE FROM %s WHERE id=$1 and project_id=$2",
			tableProjectUsers,
		),
		userID,
		projectID,
	)
	return err
}

// DropProjectUsers drops the mongo collection of this project's users
func (d *Database) DropProjectUsers(projectID string) error {
	_, err := d.db.Exec(
		fmt.Sprintf(
			"DELETE FROM %s WHERE project_id=$1",
			tableProjectUsers,
		),
		projectID,
	)
	return err
}
