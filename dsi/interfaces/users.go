package interfaces

import "github.com/machinable/machinable/dsi/models"

// UsersDatastore exposes functions to manage application users
type UsersDatastore interface {
	GetAppUserByUsername(userName string) (*models.User, error)
	GetAppUserByID(id string) (*models.User, error)
	CreateAppUser(user *models.User) error
	UpdateUserPassword(userID, passwordHash string) error
	ActivateUser(userID string, active bool) error
}
