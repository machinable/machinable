package mongo

import (
	"context"
	"time"

	"bitbucket.org/nsjostrom/machinable/dsi/models"
	"github.com/mongodb/mongo-go-driver/bson"
	"github.com/mongodb/mongo-go-driver/mongo/findopt"
)

// AddProjectLog saves a new log for a project
func (d *Datastore) AddProjectLog(project string, log *models.Log) error {
	// Get the logs collection
	col := d.db.LogDocs(project)
	_, err := col.InsertOne(
		context.Background(),
		log,
	)

	return err
}

// GetProjectLogsForLastHours retrieves all logs for a project within the last x `hours`
func (d *Datastore) GetProjectLogsForLastHours(project string, hours int) ([]*models.Log, error) {
	logs := make([]*models.Log, 0)

	// get project log collection
	collection := d.db.LogDocs(project)

	// sort by created, descending
	sortOpt := findopt.Sort(bson.NewDocument(
		bson.EC.Int32("created", -1),
	))

	// filter anything within x hours
	old := time.Now().Add(-time.Hour * time.Duration(hours))
	filterOpt := bson.NewDocument(
		bson.EC.SubDocumentFromElements("created",
			bson.EC.DateTime("$gte", old.UnixNano()/int64(time.Millisecond)),
		),
	)

	// Find logs
	cursor, err := collection.Find(
		context.Background(),
		filterOpt,
		sortOpt,
	)

	if err != nil {
		return logs, err
	}

	for cursor.Next(context.Background()) {
		var log models.Log
		err := cursor.Decode(&log)
		if err != nil {
			return logs, err
		}
		logs = append(logs, &log)
	}

	return logs, nil
}
