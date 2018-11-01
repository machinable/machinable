package handlers

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/mongodb/mongo-go-driver/bson"
	"github.com/mongodb/mongo-go-driver/bson/objectid"

	"bitbucket.org/nsjostrom/machinable/management/database"
	"bitbucket.org/nsjostrom/machinable/management/models"
	"github.com/gin-gonic/gin"
)

// UpdateProject updates the project settings, specifically the authn value
func UpdateProject(c *gin.Context) {
	projectSlug := c.Param("projectSlug")
	requestingUserID := c.MustGet("user_id").(string)

	// Create object ID from resource ID string
	userID, err := objectid.FromHex(requestingUserID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("invalid user identifier '%s': %s", requestingUserID, err.Error())})
		return
	}

	projectCollection := database.Connect().Collection(database.Projects)
	// Find project based on ID and UserID
	documentResult := projectCollection.FindOne(
		nil,
		bson.NewDocument(
			bson.EC.String("slug", projectSlug),
			bson.EC.ObjectID("user_id", userID),
		),
		nil,
	)

	if documentResult == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "no documents for collection"})
	}

	// Decode result into document
	project := &models.Project{}
	err = documentResult.Decode(project)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": fmt.Sprintf("project does not exist, '%s'", projectSlug)})
		return
	}

	var updatedProject models.ProjectBody
	c.BindJSON(&updatedProject)

	// only updated `authn` value
	_, err = projectCollection.UpdateOne(
		context.Background(),
		bson.NewDocument(
			bson.EC.String("slug", projectSlug),
			bson.EC.ObjectID("user_id", userID),
		),
		bson.NewDocument(
			bson.EC.SubDocumentFromElements("$set",
				bson.EC.Boolean("authn", updatedProject.Authn),
			),
		),
	)

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	// update project for response
	project.Authn = updatedProject.Authn

	// successful update
	c.JSON(http.StatusOK, project)
}

// CreateProject creates a new project for an application user.
func CreateProject(c *gin.Context) {
	var newProject models.ProjectBody
	requestingUserID := c.MustGet("user_id").(string)

	c.BindJSON(&newProject)
	// set user ID based on jwt
	newProject.UserID = requestingUserID

	// validate project
	err := newProject.Validate()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// get projects collection and check for duplicate slug
	pc := database.Connect().Collection(database.Projects)
	if newProject.DuplicateSlug(pc) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "project slug is already in use"})
		return
	}

	// create ObjectID from UserID string
	userObjectID, err := objectid.FromHex(newProject.UserID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// init project struc
	project := &models.Project{
		UserID:      userObjectID,
		Slug:        newProject.Slug,
		Name:        newProject.Name,
		Description: newProject.Description,
		Icon:        newProject.Icon,
		Created:     time.Now(),
		Authn:       newProject.Authn,
	}

	// save user project
	_, err = pc.InsertOne(
		context.Background(),
		project,
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// return created project to user
	c.JSON(http.StatusCreated, project)
}

// ListUserProjects returns the complete list of projects for an application user.
func ListUserProjects(c *gin.Context) {
	// grab user id from request context
	userID := c.MustGet("user_id").(string)

	// connect to project collection
	projectCollection := database.Connect().Collection(database.Projects)

	// create ObjectID from UserID string
	userObjectID, err := objectid.FromHex(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// look up projects
	cursor, err := projectCollection.Find(
		nil,
		bson.NewDocument(
			bson.EC.ObjectID("user_id", userObjectID),
		),
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	projects := make([]*models.Project, 0)
	for cursor.Next(context.Background()) {
		prj := &models.Project{}
		err := cursor.Decode(prj)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		projects = append(projects, prj)
	}
	c.JSON(http.StatusOK, gin.H{"items": projects})
}

// DeleteUserProject completely deletes an application user's project, including all related DB collections.
func DeleteUserProject(c *gin.Context) {
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
