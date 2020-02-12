package models

import "encoding/json"

// WebHook defines the structure of a project web hook
type WebHook struct {
	ID        string `json:"id"`
	ProjectID string `json:"project_id"`
	Label     string `json:"label"`
	IsEnabled bool   `json:"is_enabled"`
	Entity    string `json:"entity"`
	EntityID  string `json:"entity_id"`
	HookEvent string `json:"event"`
	Headers   []byte `json:"headers"`
	HookURL   string `json:"hook_url"`
}

// MarshalJSON custom marshaller to marshall properties to json
func (w *WebHook) MarshalJSON() ([]byte, error) {
	headers := []map[string]string{}
	err := json.Unmarshal(w.Headers, &headers)
	if err != nil {
		return nil, err
	}

	return json.Marshal(&struct {
		ID        string              `json:"id"`
		ProjectID string              `json:"project_id"`
		Label     string              `json:"label"`
		IsEnabled bool                `json:"is_enabled"`
		Entity    string              `json:"entity"`
		EntityID  string              `json:"entity_id"`
		HookEvent string              `json:"event"`
		Headers   []map[string]string `json:"headers"`
		HookURL   string              `json:"hook_url"`
	}{
		ID:        w.ID,
		ProjectID: w.ProjectID,
		Label:     w.Label,
		IsEnabled: w.IsEnabled,
		Entity:    w.Entity,
		EntityID:  w.EntityID,
		HookEvent: w.HookEvent,
		Headers:   headers,
		HookURL:   w.HookURL,
	})
}

// UnmarshalJSON is a custom unmarshaller
func (h *WebHook) UnmarshalJSON(b []byte) error {
	payload := struct {
		ID        string          `json:"id"`
		ProjectID string          `json:"project_id"`
		Label     string          `json:"label"`
		IsEnabled bool            `json:"is_enabled"`
		Entity    string          `json:"entity"`
		EntityID  string          `json:"entity_id"`
		HookEvent string          `json:"event"`
		Headers   json.RawMessage `json:"headers"`
		HookURL   string          `json:"hook_url"`
	}{}

	err := json.Unmarshal(b, &payload)

	if err != nil {
		panic(err)
	}

	h.ID = payload.ID
	h.ProjectID = payload.ProjectID
	h.Label = payload.Label
	h.IsEnabled = payload.IsEnabled
	h.Entity = payload.Entity
	h.EntityID = payload.EntityID
	h.HookEvent = payload.HookEvent
	h.Headers = payload.Headers
	h.HookURL = payload.HookURL

	return nil
}
