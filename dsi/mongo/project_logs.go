package mongo

import (
	"context"

	"github.com/anothrnick/machinable/dsi/models"
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

// CountProjectLogs returns the count of logs for a project
func (d *Datastore) CountProjectLogs(project string, filter *models.Filters) (int64, error) {
	collection := d.db.LogDocs(project)
	filterOpt, err := filtersToDocument(filter)
	if err != nil {
		return 0, err
	}
	cnt, err := collection.CountDocuments(nil, filterOpt, nil)

	return cnt, err
}

// ListProjectLogs retrieves logs based on the limit, offset, filter, and sort parameters
func (d *Datastore) ListProjectLogs(project string, limit, offset int64, filter *models.Filters, sort map[string]int) ([]*models.Log, error) {
	logs := make([]*models.Log, 0)

	// get project log collection
	collection := d.db.LogDocs(project)

	limitOpt := findopt.Limit(limit)
	offsetOpt := findopt.Skip(offset)
	sortOpt := findopt.Sort(sort)
	filterOpt, err := filtersToDocument(filter)
	if err != nil {
		return logs, err
	}

	cursor, err := collection.Find(
		context.Background(),
		filterOpt,
		limitOpt,
		offsetOpt,
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
