package projects

import (
	"net/http"

	"github.com/anothrnick/machinable/dsi/interfaces"
	"github.com/anothrnick/machinable/dsi/models"
	"github.com/gin-gonic/gin"
)

// New returns a pointer to a new `Projects`
func New(db interfaces.Datastore) *Projects {
	return &Projects{
		store: db,
	}
}

// Projects contains the datastore and any HTTP handlers needed for application projects
type Projects struct {
	store interfaces.Datastore
}

// UpdateProject updates the project settings, specifically the authn value
func (p *Projects) UpdateProject(c *gin.Context) {
	projectSlug := c.Param("projectSlug")
	userID := c.MustGet("user_id").(string)

	var updatedProject ProjectBody
	c.BindJSON(&updatedProject)

	project, err := p.store.UpdateProject(projectSlug, userID, &models.Project{Authn: updatedProject.Authn, UserRegistration: updatedProject.UserRegistration})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	project.Slug = projectSlug

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
		newProject.UserRegistration,
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
	projectSlug := c.Param("projectSlug")
	// grab user id from request context
	userID := c.MustGet("user_id").(string)

	// be sure this user owns the project
	_, err := p.store.GetProjectBySlugAndUserID(projectSlug, userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "project does not exist"})
		return
	}

	logErr := p.store.DropProjectLogs(projectSlug)
	if logErr != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "error deleting project logs"})
		return
	}
	keyErr := p.store.DropProjectKeys(projectSlug)
	if keyErr != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "error deleting project api keys"})
		return
	}
	usersErr := p.store.DropProjectUsers(projectSlug)
	if usersErr != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "error deleting project users"})
		return
	}
	sessionsErr := p.store.DropProjectSessions(projectSlug)
	if sessionsErr != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "error deleting project sessions"})
		return
	}
	resourceErr := p.store.DropProjectResources(projectSlug)
	if resourceErr != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "error deleting project resources"})
		return
	}
	projectErr := p.store.DeleteProject(projectSlug)
	if projectErr != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "error deleting project"})
		return
	}

	c.JSON(http.StatusNoContent, gin.H{})
}
