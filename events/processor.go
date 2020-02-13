package events

import "github.com/go-redis/redis"

// Event defines the event to be processed
type Event struct {
	ProjectID string `json:"project_id"`
	Entity    string `json:"entity"` // resource, json
	EntityID  string `json:"entity_id"`
	Action    string `json:"action"` // create, edit, delete
	Payload   []byte `json:"payload"`
}

// Processor process and emits events for web hooks and websockets
type Processor struct {
	cache redis.UniversalClient
}

// PushEvent processes and emits an event
func (p *Processor) PushEvent(e *Event) error {
	return nil
}
