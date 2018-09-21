package main

import (
	"net/http"

	"bitbucket.org/nsjostrom/machinable/handlers"
	"bitbucket.org/nsjostrom/machinable/middleware"
	"github.com/gin-gonic/gin"
)

func notImplemented(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{})
}

func migrate() {
}

func main() {
	router := gin.Default()
	router.Use(middleware.OpenCORSMiddleware())

	resources := router.Group("/resources")
	resources.POST("/", handlers.AddResourceDefinition)

	api := router.Group("/api")
	api.POST("/:resourcePathName", notImplemented)
	api.GET("/:resourcePathName", notImplemented)
	api.GET("/:resourcePathName/:resourceID", notImplemented)

	router.Run(":5001")
}
