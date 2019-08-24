package postgres

import "github.com/anothrnick/machinable/dsi/models"

const tableProjectUsers = "project_users"

// GetUserByUsername retrieves a project user by the user's username
func (d *Database) GetUserByUsername(project, userName string) (*models.ProjectUser, error) {
	return nil, nil
}

// GetUserByID retrieves a project user by user _id
func (d *Database) GetUserByID(project, userID string) (*models.ProjectUser, error) {
	return nil, nil
}

// CreateUser creates a new project user for the project
func (d *Database) CreateUser(project string, user *models.ProjectUser) error {
	return nil
}

// UpdateUser updates the project user's access and role
func (d *Database) UpdateUser(project, userID string, user *models.ProjectUser) error {
	return nil
}

// ListUsers returns all project users for a project
func (d *Database) ListUsers(project string) ([]*models.ProjectUser, error) {
	return nil, nil
}

// DeleteUser deletes a project user for a project based on userID
func (d *Database) DeleteUser(project, userID string) error {
	return nil
}

// DropProjectUsers drops the mongo collection of this project's users
func (d *Database) DropProjectUsers(project string) error {
	return nil
}
