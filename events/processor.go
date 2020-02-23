package events

import (
	"encoding/json"
	"log"

	"github.com/anothrnick/machinable/dsi/interfaces"
	"github.com/anothrnick/machinable/dsi/models"
	"github.com/go-redis/redis"
)

// Processor process and emits events for web hooks and websockets
type Processor struct {
	cache redis.UniversalClient
	store interfaces.ProjectHooksDatastore
}

// NewProcessor creates and returns a new instance of `Processor` with the given redis client
func NewProcessor(cache redis.UniversalClient, store interfaces.ProjectHooksDatastore) *Processor {
	return &Processor{
		cache: cache,
		store: store,
	}
}

// ProcessResults listens for web hook results on the redis queue. This function should be run as a goroutine.
func (p *Processor) ProcessResults() error {
	for {
		// endlessly read from queue
		result, err := p.cache.BLPop(0, QueueHookResults).Result()

		// exit on a read error
		if err != nil {
			log.Println(err)
			return err
		}

		// unmarshal event
		hookResult := &models.HookResult{}
		if err := json.Unmarshal([]byte(result[1]), hookResult); err != nil {
			log.Println(err)
			continue
		}

		// save hook result
		if err := p.store.AddResult(hookResult); err != nil {
			log.Println(err)
			continue
		}
	}
}

// PushEvent processes and emits an event
func (p *Processor) PushEvent(e *Event) error {
	hooks := e.Project.Hooks
	for _, hook := range hooks {
		// emit event to redis for the event action
		if hook.HookEvent == e.Action &&
			hook.Entity == e.Entity &&
			hook.EntityID == e.EntityID &&
			hook.IsEnabled == true {
			hookEvent := &HookEvent{}

			if hook.Entity == "json" {
				var payload interface{}
				json.Unmarshal(e.Payload, &payload)

				container := map[string]interface{}{
					"data": payload,
					"keys": e.Keys,
				}

				hookEvent.Hook = hook
				hookEvent.Payload = container
			} else {
				var payload interface{}
				json.Unmarshal(e.Payload, &payload)
				hookEvent.Hook = hook
				hookEvent.Payload = payload
			}
			hookEvent.EntityKey = e.EntityKey

			b, merr := json.Marshal(hookEvent)
			if merr != nil {
				log.Println(merr)
				continue
			}
			if err := p.cache.RPush(QueueHooks, b).Err(); err != nil {
				log.Println(err)
			}
		}
	}
	return nil
}
