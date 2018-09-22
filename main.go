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
	resources.GET("/", handlers.ListResourceDefinitions)
	resources.GET("/:resourceDefinitionID", handlers.GetResourceDefinition)

	api := router.Group("/api")
	api.POST("/:resourcePathName", handlers.AddObject)
	api.GET("/:resourcePathName", handlers.ListObjects)
	api.GET("/:resourcePathName/:resourceID", handlers.GetObject)

	api.POST("/:resourcePathName/:resourceID/:childResourcePathName", notImplemented)
	api.GET("/:resourcePathName/:resourceID/:childResourcePathName", notImplemented)

	router.Run(":5001")
}
