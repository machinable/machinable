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

func main() {
	router := gin.Default()
	router.Use(middleware.OpenCORSMiddleware())

	tests := router.Group("/tests")
	tests.POST("/", handlers.AddTest)
	tests.GET("/", handlers.GetTests)
	tests.DELETE("/", handlers.DeleteTests)

	collections := router.Group("/collections")
	collections.GET("/", handlers.GetCollections)
	collections.POST("/", handlers.AddCollection)
	collections.GET("/:collectionName", notImplemented)
	collections.GET("/:collectionName/*collectionKeys", notImplemented)

	resources := router.Group("/resources")
	resources.POST("/", handlers.AddResourceDefinition)
	resources.GET("/", handlers.ListResourceDefinitions)
	resources.GET("/:resourceDefinitionID", handlers.GetResourceDefinition)
	resources.DELETE("/:resourceDefinitionID", handlers.DeleteResourceDefinition)

	api := router.Group("/api")
	api.POST("/:resourcePathName", handlers.AddObject)
	api.GET("/:resourcePathName", handlers.ListObjects)
	api.GET("/:resourcePathName/:resourceID", handlers.GetObject)
	api.DELETE("/:resourcePathName/:resourceID", handlers.DeleteObject)

	//api.POST("/:resourcePathName/:resourceID/:childResourcePathName", notImplemented)
	//api.GET("/:resourcePathName/:resourceID/:childResourcePathName", notImplemented)

	router.Run(":5001")
}
