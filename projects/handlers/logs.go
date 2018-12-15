package handlers

import (
	"context"
	"net/http"
	"time"

	"bitbucket.org/nsjostrom/machinable/projects/database"
	"bitbucket.org/nsjostrom/machinable/projects/models"
	"github.com/gin-gonic/gin"
	"github.com/mongodb/mongo-go-driver/bson"
	"github.com/mongodb/mongo-go-driver/mongo/findopt"
)

// GetProjectLogs returns the list of activity logs for a project
func GetProjectLogs(c *gin.Context) {
	projectSlug := c.MustGet("project").(string)
	logs := make([]*models.Log, 0)

	collection := database.Collection(database.LogDocs(projectSlug))

	sortOpt := findopt.Sort(bson.NewDocument(
		bson.EC.Int32("created", -1),
	))

	old := time.Now().Add(-time.Hour * 24)
	filterOpt := bson.NewDocument(
		bson.EC.SubDocumentFromElements("created",
			bson.EC.DateTime("$gte", old.UnixNano()/int64(time.Millisecond)),
		),
	)
	// Find all logs
	cursor, err := collection.Find(
		context.Background(),
		filterOpt,
		sortOpt,
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	for cursor.Next(context.Background()) {
		var log models.Log
		err := cursor.Decode(&log)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		logs = append(logs, &log)
	}

	c.IndentedJSON(http.StatusOK, gin.H{"items": logs})
}
