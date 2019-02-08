package models

import "time"

// Collection is a mongo collection containing whatever data a user
// wants to save
type Collection struct {
	ID            string    `json:"id"`
	Name          string    `json:"name"`
	ParallelRead  bool      `json:"parallel_read"`
	ParallelWrite bool      `json:"parallel_write"`
	Created       time.Time `json:"created"`
	Items         int64     `json:"items"`
}
