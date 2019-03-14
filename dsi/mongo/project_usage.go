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
	Timestamp    int64   `bson:"timestamp"`     // timestamp in unix time, i.e. number of seconds elapsed since January 1, 1970 UTC
	ResponseTime float64 `bson:"response_time"` // milliseconds
}

// MResponseTimes records the response times of requests to collections and api resources over a 5 minute interval
type MResponseTimes struct {
	ID            objectid.ObjectID `bson:"_id,omitempty"`
	Type          string            `bson:"type"` // type will be response_time`
	Store         string            `bson:"store"`
	Timestamp     int64             `bson:"timestamp"`      // timestamp in unix time, i.e. number of seconds elapsed since January 1, 1970 UTC
	ResponseTimes []MResponseTiming `bson:"response_times"` // milliseconds
}

// MStatusCode records status codes of requests to collections and api resources over a 5 minute interval
type MStatusCode struct {
	ID        objectid.ObjectID `bson:"_id,omitempty"`
	Type      string            `bson:"type"` // type will be status_code
	Store     string            `bson:"store"`
	Timestamp int64             `bson:"timestamp"` // timestamp in unix time, i.e. number of seconds elapsed since January 1, 1970 UTC
	Codes     map[string]int64  `bson:"codes"`     // a map of status codes to the count
}

/* API RESOURCES */

// SaveResourceResponseTimes saves the responses times, which may overwrite
func (d *Datastore) SaveResourceResponseTimes(project string, timestamp int64, responseTimes *models.ResponseTimes) error {
	return d.saveResponseTimes(project, timestamp, responseTimes, TypeResponseTimes, ResourceStore)
}

// ListResourceResponseTimes returns a list of response times based on the filter
func (d *Datastore) ListResourceResponseTimes(project string, filter *models.Filters) ([]*models.ResponseTimes, error) {
	return d.listResponseTimes(project, filter, TypeResponseTimes, ResourceStore)
}

// SaveResourceStatusCode saves the status codes for that timestamp, may overwrite
func (d *Datastore) SaveResourceStatusCode(project string, timestamp int64, statusCode *models.StatusCode) error {
	return d.saveStatusCode(project, timestamp, statusCode, TypeStatusCode, ResourceStore)
}

// ListResourceStatusCode lists all status codes with timestamps based on the filter
func (d *Datastore) ListResourceStatusCode(project string, filter *models.Filters) ([]*models.StatusCode, error) {
	return d.listStatusCode(project, filter, TypeStatusCode, ResourceStore)
}

/* COLLECTIONS */

// SaveCollectionResponseTimes saves the responses times, which may overwrite
func (d *Datastore) SaveCollectionResponseTimes(project string, timestamp int64, responseTimes *models.ResponseTimes) error {
	return d.saveResponseTimes(project, timestamp, responseTimes, TypeResponseTimes, CollectionStore)
}

// ListCollectionResponseTimes returns a list of response times based on the filter
func (d *Datastore) ListCollectionResponseTimes(project string, filter *models.Filters) ([]*models.ResponseTimes, error) {
	return d.listResponseTimes(project, filter, TypeResponseTimes, CollectionStore)
}

// SaveCollectionStatusCode saves the status codes for that timestamp, may overwrite
func (d *Datastore) SaveCollectionStatusCode(project string, timestamp int64, statusCode *models.StatusCode) error {
	return d.saveStatusCode(project, timestamp, statusCode, TypeStatusCode, CollectionStore)
}

// ListCollectionStatusCode lists all status codes with timestamps based on the filter
func (d *Datastore) ListCollectionStatusCode(project string, filter *models.Filters) ([]*models.StatusCode, error) {
	return d.listStatusCode(project, filter, TypeStatusCode, CollectionStore)
}

func (d *Datastore) listStatusCode(project string, filter *models.Filters, typ, store string) ([]*models.StatusCode, error) {
	filterOpt, err := filtersToDocument(filter)
	if err != nil {
		return nil, err
	}

	collection := d.db.UsageDocs(project)

	filterOpt = filterOpt.Append(
		bson.EC.String("type", typ),
		bson.EC.String("store", store),
	)

	// look up the response times
	cursor, err := collection.Find(
		nil,
		filterOpt,
		nil,
	)

	if err != nil {
		return nil, err
	}

	statusCodes := make([]*models.StatusCode, 0)
	for cursor.Next(context.Background()) {
		mCode := &MStatusCode{}
		err := cursor.Decode(mCode)
		if err == nil {
			statusCodes = append(statusCodes, &models.StatusCode{
				Timestamp: mCode.Timestamp,
				Codes:     mCode.Codes,
			})
		}
	}

	return statusCodes, nil
}

func (d *Datastore) saveStatusCode(project string, timestamp int64, statusCode *models.StatusCode, typ, store string) error {
	// get a connection to the usage collection for the project
	collection := d.db.UsageDocs(project)

	mStatusCode := &MStatusCode{}
	// look up the response times
	collection.FindOne(
		nil,
		bson.NewDocument(
			bson.EC.String("type", typ),
			bson.EC.String("store", store),
			bson.EC.Int64("timestamp", timestamp),
		),
		nil,
	).Decode(mStatusCode)

	// set appropriate information if this is the first record for the aligned time
	mStatusCode.Timestamp = statusCode.Timestamp
	mStatusCode.Type = typ
	mStatusCode.Store = store

	// initialize the map if needed
	mCodes := mStatusCode.Codes
	if len(mCodes) == 0 {
		mCodes = make(map[string]int64)
	}

	// set the existing counts, update
	for code, count := range statusCode.Codes {
		if _, ok := mCodes[code]; !ok {
			mCodes[code] = 0
		}
		mCodes[code] += count
	}
	mStatusCode.Codes = mCodes

	// update record for timestamp
	_, err := collection.ReplaceOne(
		nil,
		bson.NewDocument(
			bson.EC.String("type", typ),
			bson.EC.String("store", store),
			bson.EC.Int64("timestamp", timestamp),
		),
		mStatusCode,
		replaceopt.OptUpsert(true),
	)

	return err
}

func (d *Datastore) listResponseTimes(project string, filter *models.Filters, typ, store string) ([]*models.ResponseTimes, error) {
	filterOpt, err := filtersToDocument(filter)
	if err != nil {
		return nil, err
	}

	collection := d.db.UsageDocs(project)

	filterOpt = filterOpt.Append(
		bson.EC.String("type", typ),
		bson.EC.String("store", store),
	)

	// look up the response times
	cursor, err := collection.Find(
		nil,
		filterOpt,
		nil,
	)

	if err != nil {
		return nil, err
	}

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

func (d *Datastore) saveResponseTimes(project string, timestamp int64, responseTimes *models.ResponseTimes, typ, store string) error {
	collection := d.db.UsageDocs(project)

	mResponseTimes := &MResponseTimes{}
	// look up the response times
	collection.FindOne(
		nil,
		bson.NewDocument(
			bson.EC.String("type", typ),
			bson.EC.String("store", store),
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
	mResponseTimes.Type = typ
	mResponseTimes.Store = store

	// update
	_, err := collection.ReplaceOne(
		nil,
		bson.NewDocument(
			bson.EC.String("type", typ),
			bson.EC.Int64("timestamp", timestamp),
			bson.EC.String("store", store),
		),
		mResponseTimes,
		replaceopt.OptUpsert(true),
	)

	return err
}
