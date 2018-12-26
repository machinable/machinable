package projects

import (
	"net/http"

	"bitbucket.org/nsjostrom/machinable/dsi/interfaces"
	"github.com/gin-gonic/gin"
)

// New returns a pointer to a new `Projects`
func New(db interfaces.ProjectsDatastore) *Projects {
	return &Projects{
		store: db,
	}
}

// Projects contains the datastore and any HTTP handlers needed for application projects
type Projects struct {
	store interfaces.ProjectsDatastore
}

// UpdateProject updates the project settings, specifically the authn value
func (p *Projects) UpdateProject(c *gin.Context) {
	projectSlug := c.Param("projectSlug")
	userID := c.MustGet("user_id").(string)

	var updatedProject ProjectBody
	c.BindJSON(&updatedProject)

	project, err := p.store.UpdateProjectAuthn(projectSlug, userID, updatedProject.Authn)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// successful update
	c.JSON(http.StatusOK, project)
}

// CreateProject creates a new project for an application user.
func (p *Projects) CreateProject(c *gin.Context) {
	var newProject ProjectBody
	userID := c.MustGet("user_id").(string)

	c.BindJSON(&newProject)
	// set user ID based on jwt
	newProject.UserID = userID

	// validate project
	if err := newProject.Validate(); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// check for reserved slug
	if newProject.ReservedSlug() {
		c.JSON(http.StatusBadRequest, gin.H{"error": "project slug is already in use"})
		return
	}

	// check for duplicate slug
	if _, err := p.store.GetProjectBySlug(newProject.Slug); err == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "project slug is already in use"})
		return
	}

	project, err := p.store.CreateProject(
		newProject.UserID,
		newProject.Slug,
		newProject.Name,
		newProject.Description,
		newProject.Icon,
		newProject.Authn,
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// return created project to user
	c.JSON(http.StatusCreated, project)
}

// ListUserProjects returns the complete list of projects for an application user.
func (p *Projects) ListUserProjects(c *gin.Context) {
	// grab user id from request context
	userID := c.MustGet("user_id").(string)

	projects, err := p.store.ListUserProjects(userID)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"items": projects})
}

// DeleteUserProject completely deletes an application user's project, including all related DB collections.
func (p *Projects) DeleteUserProject(c *gin.Context) {
	// TODO
	// delete from projects
	// delete project collections
	// delete project resources
	// delete project resource data
	// delete project users
	// delete project keys
	// delete project logs
	// delete project usage
}
