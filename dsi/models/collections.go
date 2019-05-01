package models

import (
	"errors"
	"time"

	"github.com/anothrnick/machinable/dsi"
)

// Collection is a mongo collection containing whatever data a user
// wants to save
type Collection struct {
	ID            string    `json:"id"`
	Name          string    `json:"name"`
	ParallelRead  bool      `json:"parallel_read"`
	ParallelWrite bool      `json:"parallel_write"`
	Create        bool      `json:"create"`
	Read          bool      `json:"read"`
	Update        bool      `json:"update"`
	Delete        bool      `json:"delete"`
	Created       time.Time `json:"created"`
	Items         int64     `json:"items"`
}

// Validate validates the collection fields
func (c *Collection) Validate() error {
	if c.Name == "" {
		return errors.New("collection name cannot be empty")
	} else if !dsi.ValidPathFormat.MatchString(c.Name) {
		return errors.New("invalid collection name: only alphanumeric, dashes, and underscores allowed")
	}
	return nil
}
