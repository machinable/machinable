package middleware

import (
	"bytes"
	"log"
	"time"

	"github.com/anothrnick/machinable/dsi/interfaces"
	"github.com/anothrnick/machinable/dsi/models"
	"github.com/gin-gonic/gin"
)

type logWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (w logWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}

// JSONStatsMiddleware logs json stats for reporting
func JSONStatsMiddleware(store interfaces.Datastore) gin.HandlerFunc {
	return loggingMiddleware(store, models.EndpointJSON)
}

// ResourceStatsMiddleware logs resource stats and logging for reporting
func ResourceStatsMiddleware(store interfaces.Datastore) gin.HandlerFunc {
	return loggingMiddleware(store, models.EndpointResource)
}

func loggingMiddleware(store interfaces.Datastore, endpointType string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// inject custom writer
		lw := &logWriter{body: bytes.NewBufferString(""), ResponseWriter: c.Writer}
		c.Writer = lw

		// response time
		requestStart := time.Now()
		// get aligned time by 5 minute interval
		alignedStart := AlignTime(requestStart, 5)

		// continue handler chain
		c.Next()

		// response time in ms
		responseTime := time.Now().Sub(requestStart).Seconds() * 1000

		// get status code
		statusCode := c.Writer.Status()

		verb := c.Request.Method
		path := c.Request.URL.Path

		projectID := c.GetString("projectId")
		authType := c.GetString("authType")
		authString := c.GetString("authString")
		authID := c.GetString("authID")

		if authType == "" {
			authString = "anonymous"
			authType = "anonymous"
		}

		// save project log
		plog := &models.Log{
			EndpointType:   endpointType,
			Verb:           verb,
			Path:           path,
			StatusCode:     statusCode,
			Created:        time.Now().Unix(),
			AlignedCreated: alignedStart.Unix(),
			ResponseTime:   int64(responseTime),
			Initiator:      authString,
			InitiatorType:  authType,
			InitiatorID:    authID,
		}

		// save in go routine, do not block request
		go func(projectID string, plog *models.Log) {
			// save the log
			err := store.AddProjectLog(projectID, plog)

			if err != nil {
				log.Println("an error occured trying to save the log")
				log.Println(err.Error())
			}
		}(projectID, plog)
	}
}

// AlignTime returns the aligned `time.Time` based on the `unaligned` parameter and the `interval` to align with (in minutes)
func AlignTime(unaligned time.Time, interval int) time.Time {
	timeToAlign := unaligned.Truncate(time.Minute)
	timeOffset := (timeToAlign.Minute() % interval)
	timeAligned := timeToAlign.Add(-time.Duration(timeOffset) * time.Minute)
	return timeAligned
}
