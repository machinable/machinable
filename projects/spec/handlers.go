package spec

import (
	"net/http"

	"github.com/anothrnick/machinable/dsi/interfaces"
	"github.com/gin-gonic/gin"
)

// New returns a pointer to a new `Users` struct
func New(db interfaces.ResourcesDatastore) *Spec {
	return &Spec{
		store: db,
	}
}

// Spec wraps the datastore and any HTTP handlers for project openapi spec
type Spec struct {
	store interfaces.ResourcesDatastore
}

// GetSpec retrieves the openapi spec for the project
func (s *Spec) GetSpec(c *gin.Context) {
	spec := baseSpec()

	c.IndentedJSON(http.StatusOK, gin.H{"spec": spec})
}
