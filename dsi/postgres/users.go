package postgres

import "github.com/anothrnick/machinable/dsi/models"

const tableAppUsers = "app_users"

// GetAppUserByUsername attempts to find a user by username, if the user does not exist the user will be nil, error will be !nil
func (d *Database) GetAppUserByUsername(userName string) (*models.User, error) {
	return nil, nil
}

// GetAppUserByID attempts to find a user by ID, if the user does not exist the user will be nil, error will be !nil
func (d *Database) GetAppUserByID(id string) (*models.User, error) {
	return nil, nil
}

// CreateAppUser saves a new application user
func (d *Database) CreateAppUser(user *models.User) error {
	return nil
}

// UpdateUserPassword updates the user's password
func (d *Database) UpdateUserPassword(userID, passwordHash string) error {
	return nil
}
