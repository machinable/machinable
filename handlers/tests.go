package handlers

import (
	"context"
	"net/http"

	"bitbucket.org/nsjostrom/machinable/database"
	"github.com/gin-gonic/gin"
	"github.com/mongodb/mongo-go-driver/bson"
	"github.com/mongodb/mongo-go-driver/bson/objectid"
)

type errorItem struct {
	ID objectid.ObjectID `bson:"_id"`
}

// AddTest adds a new document
func AddTest(c *gin.Context) {
	bdoc := make(map[string]interface{})

	c.BindJSON(&bdoc)

	// Get a connection and insert the new document
	collection := database.Connect().Collection(database.Tests)
	_, err := collection.InsertOne(
		context.Background(),
		bdoc,
	)

	// TODO
	// Load result id and try to decode. If an error occurs, delete the document and return the error message with a 400 status code

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, bdoc)
}

// GetTests returns the full list of documents
func GetTests(c *gin.Context) {
	documents := make([]map[string]interface{}, 0)

	collection := database.Connect().Collection(database.Tests)

	cursor, err := collection.Find(
		context.Background(),
		bson.NewDocument(),
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	//doc := bson.NewDocument()
	for cursor.Next(context.Background()) {
		//doc.Reset()
		doc := make(map[string]interface{})
		err := cursor.Decode(doc)
		if err != nil {
			var errid string
			item := &errorItem{}
			anotherErr := cursor.Decode(item)
			if anotherErr == nil {
				errid = item.ID.Hex()
			}
			documents = append(documents, map[string]interface{}{
				"_id":   errid,
				"error": err.Error(),
			})
		} else {
			// get stringified version of the ID
			objectID, ok := doc["_id"].(objectid.ObjectID)
			if ok {
				doc["_id"] = objectID.Hex()
			}

			documents = append(documents, doc)
		}
	}
	c.JSON(http.StatusOK, gin.H{"items": documents})
}

func DeleteTests(c *gin.Context) {

}
