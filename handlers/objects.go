package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func AddObject(c *gin.Context) {
	//resourcePathName := c.Param("resourcePathName")
	c.JSON(http.StatusNoContent, gin.H{})
}

func ListObjects(c *gin.Context) {
	//resourcePathName := c.Param("resourcePathName")
	c.JSON(http.StatusNoContent, gin.H{})
}

func GetObject(c *gin.Context) {
	//resourcePathName := c.Param("resourcePathName")
	//resourceID := c.Param("resourceID")
	c.JSON(http.StatusNoContent, gin.H{})
}
