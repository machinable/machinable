package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func LogActivity() error {
	return nil
}

// GetProjectLogs returns the list of activity logs for a project
func GetProjectLogs(c *gin.Context) {

	c.JSON(http.StatusOK, gin.H{})
}
