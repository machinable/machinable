package interfaces

import "github.com/anothrnick/machinable/dsi/models"

// ProjectsDatastore exposes functions for projects
type ProjectsDatastore interface {
	UpdateProject(slug, userID string, project *models.Project) (*models.Project, error)
	UpdateProjectAuthn(slug, userID string, authn bool) (*models.Project, error)
	UpdateProjectUserRegistration(slug, userID string, registration bool) (*models.Project, error)
	CreateProject(userID, slug, name, description, icon string, authn bool) (*models.Project, error)
	ListUserProjects(userID string) ([]*models.Project, error)
	GetProjectBySlug(slug string) (*models.Project, error)
	GetProjectBySlugAndUserID(slug, userID string) (*models.Project, error)
	DeleteProject(slug string) error
}
