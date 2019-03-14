package middleware

import (
	"bytes"
	"log"
	"strconv"
	"time"

	"github.com/anothrnick/machinable/dsi/interfaces"
	"github.com/anothrnick/machinable/dsi/models"
	"github.com/gin-gonic/gin"
)

// CollectionStatsMiddleware logs collection stats for reporting
func CollectionStatsMiddleware(store interfaces.Datastore) gin.HandlerFunc {
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

		projectSlug := c.GetString("project")

		// save response time
		err := store.SaveCollectionResponseTimes(
			projectSlug,
			alignedStart.Unix(),
			&models.ResponseTimes{
				Timestamp: alignedStart.Unix(),
				ResponseTimes: []models.ResponseTiming{
					{
						Timestamp:    requestStart.Unix(),
						ResponseTime: responseTime,
					},
				},
			},
		)
		if err != nil {
			log.Println("an error occured trying to save the response time")
			log.Println(err.Error())
		}

		// save status code
		err = store.SaveCollectionStatusCode(
			projectSlug,
			alignedStart.Unix(),
			&models.StatusCode{
				Timestamp: alignedStart.Unix(),
				Codes: map[string]int64{
					strconv.Itoa(statusCode): 1,
				},
			},
		)
		if err != nil {
			log.Println("an error occured trying to save the status code")
			log.Println(err.Error())
		}
	}
}

// ResourceStatsMiddleware logs resource stats for reporting
func ResourceStatsMiddleware(store interfaces.Datastore) gin.HandlerFunc {
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

		projectSlug := c.GetString("project")

		// save response time
		err := store.SaveResourceResponseTimes(
			projectSlug,
			alignedStart.Unix(),
			&models.ResponseTimes{
				Timestamp: alignedStart.Unix(),
				ResponseTimes: []models.ResponseTiming{
					{
						Timestamp:    requestStart.Unix(),
						ResponseTime: responseTime,
					},
				},
			},
		)
		if err != nil {
			log.Println("an error occured trying to save the response time")
			log.Println(err.Error())
		}

		// save status code
		err = store.SaveResourceStatusCode(
			projectSlug,
			alignedStart.Unix(),
			&models.StatusCode{
				Timestamp: alignedStart.Unix(),
				Codes: map[string]int64{
					strconv.Itoa(statusCode): 1,
				},
			},
		)
		if err != nil {
			log.Println("an error occured trying to save the status code")
			log.Println(err.Error())
		}
	}
}

// AlignTime returns the aligned `time.Time` based on the `unaligned` parameter and the `interval` to align with (in minutes)
func AlignTime(unaligned time.Time, interval int) time.Time {
	timeToAlign := unaligned.Truncate(time.Minute)
	timeOffset := (timeToAlign.Minute() % interval)
	timeAligned := timeToAlign.Add(-time.Duration(timeOffset) * time.Minute)
	return timeAligned
}
