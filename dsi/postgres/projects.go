package postgres

import (
	"fmt"

	"github.com/anothrnick/machinable/dsi/models"
)

const tableAppProjects = "app_projects"

// UpdateProject updates the project's name, description, icon, and user_registration
func (d *Database) UpdateProject(slug, userID string, project *models.Project) (*models.Project, error) {
	_, err := d.db.Exec(
		fmt.Sprintf(
			"UPDATE %s SET name=$1, description=$2, icon=$3, user_registration=$4  WHERE slug=$5 and user_id=$7",
			tableAppProjects,
		),
		project.Name,
		project.Description,
		project.Icon,
		project.UserRegistration,
		slug,
		userID,
	)

	return project, err
}

// UpdateProjectUserRegistration updates the project authentication policy
func (d *Database) UpdateProjectUserRegistration(slug, userID string, registration bool) (*models.Project, error) {
	_, err := d.db.Exec(
		fmt.Sprintf(
			"UPDATE %s SET user_registration=$1 WHERE slug=$2 and user_id=$3",
			tableAppProjects,
		),
		registration,
		slug,
		userID,
	)

	return nil, err
}

// CreateProject creates a new project for a user
func (d *Database) CreateProject(userID, slug, name, description, icon string, authn bool, register bool) (*models.Project, error) {
	project := models.Project{
		UserID:           userID,
		Slug:             slug,
		Name:             name,
		Description:      description,
		Icon:             icon,
		UserRegistration: register,
	}
	err := d.db.QueryRow(
		fmt.Sprintf(
			"INSERT INTO %s (user_id, slug, name, description, icon, user_registration) VALUES ($1, $2, $3, $4, $5, $6) RETURNING id",
			tableAppProjects,
		),
		project.UserID,
		project.Slug,
		project.Name,
		project.Description,
		project.Icon,
		project.UserRegistration,
	).Scan(&project.ID)

	return &project, err
}

// ListUserProjects retrieves all projects for a user
func (d *Database) ListUserProjects(userID string) ([]*models.Project, error) {
	rows, err := d.db.Query(
		fmt.Sprintf(
			"SELECT id, user_id, slug, name, description, icon, user_registration, created FROM %s WHERE user_id=$1",
			tableAppProjects,
		),
		userID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	projects := make([]*models.Project, 0)
	for rows.Next() {
		project := models.Project{}
		err = rows.Scan(
			&project.ID,
			&project.UserID,
			&project.Slug,
			&project.Name,
			&project.Description,
			&project.Icon,
			&project.UserRegistration,
			&project.Created,
		)
		if err != nil {
			return nil, err
		}

		projects = append(projects, &project)
	}

	return projects, rows.Err()
}

// GetProjectBySlug retrieves a project by slug
func (d *Database) GetProjectBySlug(slug string) (*models.Project, error) {
	project := models.Project{}

	err := d.db.QueryRow(
		fmt.Sprintf(
			"SELECT id, user_id, slug, name, description, icon, user_registration, created FROM %s WHERE slug=$1",
			tableAppProjects,
		),
		slug,
	).Scan(
		&project.ID,
		&project.UserID,
		&project.Slug,
		&project.Name,
		&project.Description,
		&project.Icon,
		&project.UserRegistration,
		&project.Created,
	)
	if err != nil {
		return nil, err
	}

	return &project, err
}

// GetProjectBySlugAndUserID retrieves a project by slug for a given user ID
func (d *Database) GetProjectBySlugAndUserID(slug, userID string) (*models.Project, error) {
	project := models.Project{}

	err := d.db.QueryRow(
		fmt.Sprintf(
			"SELECT id, user_id, slug, name, description, icon, user_registration, created FROM %s WHERE slug=$1 and user_id=$2",
			tableAppProjects,
		),
		slug,
		userID,
	).Scan(
		&project.ID,
		&project.UserID,
		&project.Slug,
		&project.Name,
		&project.Description,
		&project.Icon,
		&project.UserRegistration,
		&project.Created,
	)
	if err != nil {
		return nil, err
	}

	return &project, err
}

// DeleteProject permanently removes a project based on project slug
func (d *Database) DeleteProject(slug string) error {
	_, err := d.db.Exec(
		fmt.Sprintf(
			"DELETE FROM %s WHERE slug=$1",
			tableAppProjects,
		),
		slug,
	)
	return err
}
