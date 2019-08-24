package postgres

import "github.com/anothrnick/machinable/dsi/models"

const tableAppProjects = "app_projects"

// UpdateProject updates the project
func (d *Database) UpdateProject(slug, userID string, project *models.Project) (*models.Project, error) {
	return nil, nil
}

// UpdateProjectAuthn updates the project authentication policy
func (d *Database) UpdateProjectAuthn(slug, userID string, authn bool) (*models.Project, error) {
	return nil, nil
}

// UpdateProjectUserRegistration updates the project authentication policy
func (d *Database) UpdateProjectUserRegistration(slug, userID string, registration bool) (*models.Project, error) {
	return nil, nil
}

// CreateProject creates a new project for a user
func (d *Database) CreateProject(userID, slug, name, description, icon string, authn bool, register bool) (*models.Project, error) {
	return nil, nil
}

// ListUserProjects retrieves all projects for a user
func (d *Database) ListUserProjects(userID string) ([]*models.Project, error) {
	return nil, nil
}

// GetProjectBySlug retrieves a project by slug
func (d *Database) GetProjectBySlug(slug string) (*models.Project, error) {
	return nil, nil
}

// GetProjectBySlugAndUserID retrieves a project by slug for a given user ID
func (d *Database) GetProjectBySlugAndUserID(slug, userID string) (*models.Project, error) {
	return nil, nil
}

// DeleteProject permanently removes a project based on project slug
func (d *Database) DeleteProject(slug string) error {
	return nil
}
