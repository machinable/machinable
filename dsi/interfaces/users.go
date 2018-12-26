package interfaces

import "bitbucket.org/nsjostrom/machinable/dsi/models"

// UsersDatastore exposes functions to manage application users
type UsersDatastore interface {
	GetAppUserByUsername(userName string) (*models.User, error)
	GetAppUserByID(id string) (*models.User, error)
	CreateAppUser(user *models.User) error
}
