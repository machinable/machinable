package interfaces

import "github.com/anothrnick/machinable/dsi/models"

// ProjectUsersDatastore exposes functions to manage project users
type ProjectUsersDatastore interface {
	GetUserByUsername(project, userName string) (*models.ProjectUser, error)
	CreateUser(project string, user *models.ProjectUser) error
	ListUsers(project string) ([]*models.ProjectUser, error)
	DeleteUser(project, userID string) error
	DropProjectUsers(project string) error
}
