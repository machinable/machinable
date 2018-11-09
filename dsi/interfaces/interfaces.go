package interfaces

import "bitbucket.org/nsjostrom/machinable/dsi/models"

// ProjectDefinitionDocuments provides an interface to the collection of defined api objects
type ProjectDefinitionDocuments interface {
	AddDocument(project, path string, fields map[string]interface{}) (string, error)
	ListDocuments(project, path string, limit, offset int, filter map[string]interface{}) ([]map[string]interface{}, error)
	GetDocument(project, path, documentID string) (map[string]interface{}, error)
	DeleteDocument(project, path, documentID string) error
	DropAll(project, path string) error
}

// ProjectDefinitions provides an interface to the collection of definitions
type ProjectDefinitions interface {
	AddDefinition(project string, def *models.ResourceDefinition) (string, error)
	ListDefinitions(project string) ([]*models.ResourceDefinition, error)
	GetDefinition(project, definitionID string) (*models.ResourceDefinition, error)
	DeleteDefinition(project, definitionID string) error
}

// ProjectCollections provides an interface to the project collections. The collections are just references to the documents by the user provided path name.
type ProjectCollections interface {
	CreateCollection(project, name string) error
	GetCollection(project, name string) (string, error)
	GetCollections(project string) ([]string, error)
	DeleteCollection(project, name string) error
}

// ProjectCollectionDocuments provides an interface to the project collection documents. These documents can have any structure.
type ProjectCollectionDocuments interface {
	AddDocument(project, collectionName string, document map[string]interface{}) (map[string]interface{}, error)
	UpdateDocument(project, collectionName, documentID string, updatedDocumet map[string]interface{}) error
	GetDocuments(project, collectionName string, limit, offset int, filter map[string]interface{}) ([]map[string]interface{}, error)
	GetDocument(project, collectionName, documentID string) (map[string]interface{}, error)
	DropAll(project, collectionName string) error
}
