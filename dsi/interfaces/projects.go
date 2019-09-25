package interfaces

import "github.com/anothrnick/machinable/dsi/models"

// ProjectsDatastore exposes functions for projects
type ProjectsDatastore interface {
	UpdateProject(projectID, userID string, project *models.Project) (*models.Project, error)
	UpdateProjectUserRegistration(projectID, userID string, registration bool) (*models.Project, error)
	CreateProject(userID, slug, name, description, icon string, authn bool, register bool) (*models.Project, error)
	ListUserProjects(projectID string) ([]*models.Project, error)
	GetProjectBySlug(slug string) (*models.Project, error)
	GetProjectBySlugAndUserID(slug, userID string) (*models.Project, error)
	DeleteProject(projectID string) error
}
