package models

// WebHook defines the structure of a project web hook
type WebHook struct {
	ID        string              `json:"id"`
	ProjectID string              `json:"project_id"`
	Label     string              `json:"label"`
	IsEnabled bool                `json:"is_enabled"`
	Entity    string              `json:"entity"`
	EntityID  string              `json:"entity_id"`
	HookEvent string              `json:"event"`
	Headers   []map[string]string `json:"headers"`
	HookURL   string              `json:"hook_url"`
}
