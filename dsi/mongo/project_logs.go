package mongo

import (
	"context"

	"github.com/anothrnick/machinable/dsi/models"
	"github.com/mongodb/mongo-go-driver/mongo/findopt"
)

type MongoLog struct {
	Event         string `bson:"event"`
	StatusCode    int    `bson:"status_code"`
	Created       int64  `bson:"created"`
	Initiator     string `bson:"initiator"`
	InitiatorType string `bson:"initiator_type"`
	InitiatorID   string `bson:"initiator_id"`
	TargetID      string `bson:"target_id"`
}

// AddProjectLog saves a new log for a project
func (d *Datastore) AddProjectLog(project string, log *models.Log) error {
	mongoLog := &MongoLog{
		Event:         log.Event,
		StatusCode:    log.StatusCode,
		Created:       log.Created,
		Initiator:     log.Initiator,
		InitiatorType: log.InitiatorType,
		InitiatorID:   log.InitiatorID,
		TargetID:      log.TargetID,
	}

	// Get the logs collection
	col := d.db.LogDocs(project)
	_, err := col.InsertOne(
		context.Background(),
		mongoLog,
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
		var log MongoLog
		err := cursor.Decode(&log)
		if err != nil {
			return logs, err
		}
		logs = append(logs, &models.Log{
			Event:         log.Event,
			StatusCode:    log.StatusCode,
			Created:       log.Created,
			Initiator:     log.Initiator,
			InitiatorType: log.InitiatorType,
			InitiatorID:   log.InitiatorID,
			TargetID:      log.TargetID,
		})
	}

	return logs, nil
}
