package models

import (
	"encoding/json"
	"errors"
	"net/url"
)

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

// validURL parses the string as a url and verifies it is valid
func validURL(str string) bool {
	u, err := url.Parse(str)
	return err == nil && u.Scheme != "" && u.Host != ""
}

// Validate performs validation on the WebHook struct, returning an error if it is invalid
func (w *WebHook) Validate() error {
	if w.ProjectID == "" {
		return errors.New("invalid project id")
	} else if w.Label == "" {
		return errors.New("label can not be empty")
	} else if w.Entity == "" {
		return errors.New("entity can not be empty")
	} else if w.EntityID == "" {
		return errors.New("entity can not be empty")
	} else if w.HookEvent == "" {
		return errors.New("hook event can not be empty")
	} else if w.HookURL == "" {
		return errors.New("hook URL can not be empty")
	} else if !validURL(w.HookURL) {
		return errors.New("invalid hook URL")
	}

	return nil
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

// UnmarshalJSON is a custom unmarshaller, specificall for the `headers`
func (w *WebHook) UnmarshalJSON(b []byte) error {
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

	w.ID = payload.ID
	w.ProjectID = payload.ProjectID
	w.Label = payload.Label
	w.IsEnabled = payload.IsEnabled
	w.Entity = payload.Entity
	w.EntityID = payload.EntityID
	w.HookEvent = payload.HookEvent
	w.Headers = payload.Headers
	w.HookURL = payload.HookURL

	return nil
}
