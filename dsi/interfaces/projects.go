package interfaces

import "bitbucket.org/nsjostrom/machinable/dsi/models"

// ProjectsDatastore exposes functions for projects
type ProjectsDatastore interface {
	UpdateProjectAuthn(slug, userID string, authn bool) (*models.Project, error)
	CreateProject(userID, slug, name, description, icon string, authn bool) (*models.Project, error)
	ListUserProjects(userID string) ([]*models.Project, error)
	GetProjectBySlug(slug string) (*models.Project, error)
	GetProjectBySlugAndUserID(slug, userID string) (*models.Project, error)
	DeleteProject(slug string) error
}
