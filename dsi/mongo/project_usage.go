package mongo

import (
	"context"

	"github.com/anothrnick/machinable/dsi/models"
	"github.com/mongodb/mongo-go-driver/bson"
	"github.com/mongodb/mongo-go-driver/bson/objectid"
	"github.com/mongodb/mongo-go-driver/mongo/replaceopt"
)

const (
	CollectionStore = "collection"
	ResourceStore   = "resource"

	TypeResponseTimes = "response_time"
	TypeStatusCode    = "status_code"
)

// MResponseTiming is a single request
type MResponseTiming struct {
	Timestamp    int64 `bson:"timestamp"`     // timestamp in unix time, i.e. number of seconds elapsed since January 1, 1970 UTC
	ResponseTime int64 `bson:"response_time"` // milliseconds
}

// MResponseTimes records the response times of requests to collections and api resources over a 5 minute interval
type MResponseTimes struct {
	ID            objectid.ObjectID `bson:"_id"`
	Type          string            `bson:"type"`           // type will be response_time`
	Timestamp     int64             `bson:"timestamp"`      // timestamp in unix time, i.e. number of seconds elapsed since January 1, 1970 UTC
	ResponseTimes []MResponseTiming `bson:"response_times"` // milliseconds
}

// MStatusCode records status codes of requests to collections and api resources over a 5 minute interval
type MStatusCode struct {
	ID        objectid.ObjectID `bson:"_id"`
	Type      string            `bson:"type"`      // type will be status_code
	Timestamp int64             `bson:"timestamp"` // timestamp in unix time, i.e. number of seconds elapsed since January 1, 1970 UTC
	Codes     map[string]int64  `bson:"codes"`     // a map of status codes to the count
}

// SaveResponseTimes saves the responses times, which may overwrite
func (d *Datastore) SaveResponseTimes(project string, timestamp int64, responseTimes *models.ResponseTimes) error {
	collection := d.db.UsageDocs(project)

	mResponseTimes := &MResponseTimes{}
	// look up the response times
	collection.FindOne(
		nil,
		bson.NewDocument(
			bson.EC.String("type", CollectionStore),
			bson.EC.Int64("timestamp", timestamp),
		),
		nil,
	).Decode(mResponseTimes)

	// convert
	translatedTimes := mResponseTimes.ResponseTimes

	if responseTimes.ResponseTimes != nil {
		for _, t := range responseTimes.ResponseTimes {
			translatedTimes = append(translatedTimes, MResponseTiming{
				Timestamp:    t.Timestamp,
				ResponseTime: t.ResponseTime,
			})
		}
	}

	mResponseTimes.Timestamp = responseTimes.Timestamp
	mResponseTimes.ResponseTimes = translatedTimes

	// update
	_, err := collection.ReplaceOne(
		nil,
		bson.NewDocument(
			bson.EC.String("type", CollectionStore),
			bson.EC.Int64("timestamp", timestamp),
		),
		mResponseTimes,
		replaceopt.OptUpsert(true),
	)

	return err
}

// ListResponseTimes returns a list of response times based on the filter
func (d *Datastore) ListResponseTimes(project string) ([]*models.ResponseTimes, error) {
	collection := d.db.UsageDocs(project)

	// look up the response times
	cursor, err := collection.Find(
		nil,
		bson.NewDocument(
			bson.EC.String("type", CollectionStore),
		),
		nil,
	)

	responseTimes := make([]*models.ResponseTimes, 0)
	for cursor.Next(context.Background()) {
		mResponseTimes := &MResponseTimes{}
		err := cursor.Decode(mResponseTimes)
		if err == nil {
			// convert
			translatedTimes := []models.ResponseTiming{}

			if mResponseTimes.ResponseTimes != nil {
				for _, t := range mResponseTimes.ResponseTimes {
					translatedTimes = append(translatedTimes, models.ResponseTiming{
						Timestamp:    t.Timestamp,
						ResponseTime: t.ResponseTime,
					})
				}
			}
			responseTimes = append(responseTimes, &models.ResponseTimes{
				Timestamp:     mResponseTimes.Timestamp,
				ResponseTimes: translatedTimes,
			})
		}
	}

	return responseTimes, nil
}

func (d *Datastore) SaveStatusCode(project string, timestamp int64, statusCode *models.StatusCode) error {

}

func (d *Datastore) ListStatusCode(project string) ([]*models.StatusCode, error) {

}
