package events

import "github.com/machinable/machinable/dsi/models"

const (
	// QueueHooks is the redis queue for web hooks
	QueueHooks = "hook_queue"
	// QueueHookResults is the redis queue for web hook result messages
	QueueHookResults = "hook_result_queue"
	// QueueEmailNotifications process email send
	QueueEmailNotifications = "email_notifications_queue"
)

// Event defines the event(s) to be processed
type Event struct {
	Project   *models.ProjectDetail `json:"project"`
	Entity    string                `json:"entity"` // resource, json
	EntityKey string                `json:"entity_key"`
	EntityID  string                `json:"entity_id"`
	Action    string                `json:"action"` // create, edit, delete
	Keys      []string              `json:"keys"`
	Payload   []byte                `json:"payload"`
}

// HookEvent describes a single web hook event
type HookEvent struct {
	Hook      *models.WebHook `json:"hook"`
	EntityKey string          `json:"entity_key"`
	Payload   interface{}     `json:"payload"`
}

// Notification contains the information in the queue for an email notification
type Notification struct {
	Template         string            `json:"template"`
	Subject          string            `json:"subject"`
	ReceiverName     string            `json:"receiver_name"`
	ReceiverEmail    string            `json:"receiver_email"`
	PlainTextContent string            `json:"plain_text_content"` // plain text in the case email html can't render
	Data             map[string]string `json:"data"`               // data for template
	Meta             map[string]string `json:"meta"`               // meta data for logging
}
