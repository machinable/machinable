package events

import (
	"log"

	"github.com/anothrnick/machinable/dsi/models"
	"github.com/go-redis/redis"
)

const (
	// WebhookQueue is the queue name for web hook callbacks
	WebhookQueue = "hook_queue"
)

// Event defines the event to be processed
type Event struct {
	Project  *models.ProjectDetail `json:"project"`
	Entity   string                `json:"entity"` // resource, json
	EntityID string                `json:"entity_id"`
	Action   string                `json:"action"` // create, edit, delete
	Payload  []byte                `json:"payload"`
}

// Processor process and emits events for web hooks and websockets
type Processor struct {
	cache redis.UniversalClient
}

// PushEvent processes and emits an event
func (p *Processor) PushEvent(e *Event) error {
	hooks := e.Project.Hooks
	for _, hook := range hooks {
		// emit event to redis for the event action
		if hook.HookEvent == e.Action {
			if err := p.cache.RPush(WebhookQueue, e.Entity, e.EntityID, e.Action, e.Payload).Err(); err != nil {
				log.Println(err)
			}
		}
	}
	return nil
}
