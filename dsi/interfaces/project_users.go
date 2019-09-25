package interfaces

import "github.com/anothrnick/machinable/dsi/models"

// ProjectUsersDatastore exposes functions to manage project users
type ProjectUsersDatastore interface {
	GetUserByUsername(projectID, userName string) (*models.ProjectUser, error)
	GetUserByID(projectID, userID string) (*models.ProjectUser, error)
	CreateUser(projectID string, user *models.ProjectUser) error
	UpdateUser(projectID, userID string, user *models.ProjectUser) error
	ListUsers(projectID string) ([]*models.ProjectUser, error)
	DeleteUser(projectID, userID string) error
	DropProjectUsers(projectID string) error
}
