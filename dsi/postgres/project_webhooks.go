package postgres

import (
	"fmt"
	"time"

	"github.com/anothrnick/machinable/dsi/errors"
	"github.com/anothrnick/machinable/dsi/models"
)

const tableProjectWebhookResults = "project_webhook_results"
const tableProjectWebHooks = "project_webhooks"

// AddResult creates a new webhook result
func (d *Database) AddResult(result *models.HookResult) *errors.DatastoreError {
	_, err := d.db.Exec(
		fmt.Sprintf(
			"INSERT INTO %s (project_id, webhook_id, status_code, response_time, error_message, created) VALUES ($1, $2, $3, $4, $5, $6)",
			tableProjectWebhookResults,
		),
		result.ProjectID,
		result.WebHookID,
		result.StatusCode,
		result.ResponseTime,
		result.ErrorMessage,
		time.Now(),
	)

	return errors.New(errors.UnknownError, err)
}

// ListResults lists all webhook results for a web hook
func (d *Database) ListResults(projectID, hookID string) ([]*models.HookResult, *errors.DatastoreError) {
	rows, err := d.db.Query(
		fmt.Sprintf(
			"SELECT project_id, webhook_id, status_code, response_time, error_message, created FROM %s WHERE project_id=$1 AND webhook_id=$2 AND created >= now()-'1 hour'::interval ORDER BY created DESC",
			tableProjectWebhookResults,
		),
		projectID,
		hookID,
	)
	if err != nil {
		return nil, errors.New(errors.UnknownError, err)
	}
	defer rows.Close()

	results := make([]*models.HookResult, 0)
	for rows.Next() {
		result := models.HookResult{}
		err = rows.Scan(
			&result.ProjectID,
			&result.WebHookID,
			&result.StatusCode,
			&result.ResponseTime,
			&result.ErrorMessage,
			&result.Created,
		)
		if err != nil {
			return nil, errors.New(errors.UnknownError, err)
		}

		results = append(results, &result)
	}

	return results, nil
}

// AddHook saves a new WebHook to the datastore
func (d *Database) AddHook(projectID string, hook *models.WebHook) *errors.DatastoreError {
	_, err := d.db.Exec(
		fmt.Sprintf(
			"INSERT INTO %s (project_id, label, isenabled, entity, entity_id, hook_event, headers, hook_url) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)",
			tableProjectWebHooks,
		),
		projectID,
		hook.Label,
		hook.IsEnabled,
		hook.Entity,
		hook.EntityID,
		hook.HookEvent,
		hook.Headers,
		hook.HookURL,
	)

	return errors.New(errors.UnknownError, err)
}

// ListHooks retrieves all WebHooks for a project
func (d *Database) ListHooks(projectID string) ([]*models.WebHook, *errors.DatastoreError) {
	rows, err := d.db.Query(
		fmt.Sprintf(
			"SELECT id, project_id, label, isenabled, entity, entity_id, hook_event, headers, hook_url FROM %s WHERE project_id=$1",
			tableProjectWebHooks,
		),
		projectID,
	)
	if err != nil {
		return nil, errors.New(errors.UnknownError, err)
	}
	defer rows.Close()

	hooks := make([]*models.WebHook, 0)
	for rows.Next() {
		hook := models.WebHook{}
		err = rows.Scan(
			&hook.ID,
			&hook.ProjectID,
			&hook.Label,
			&hook.IsEnabled,
			&hook.Entity,
			&hook.EntityID,
			&hook.HookEvent,
			&hook.Headers,
			&hook.HookURL,
		)
		if err != nil {
			return nil, errors.New(errors.UnknownError, err)
		}

		hooks = append(hooks, &hook)
	}

	return hooks, nil
}

// GetHook retrieves a single hook by project and hook ID, if it exists
func (d *Database) GetHook(projectID, hookID string) (*models.WebHook, *errors.DatastoreError) {
	hook := models.WebHook{}
	err := d.db.QueryRow(
		fmt.Sprintf(
			"SELECT id, project_id, label, isenabled, entity, entity_id, hook_event, headers, hook_url FROM %s WHERE project_id=$1 AND id=$2",
			tableProjectWebHooks,
		),
		projectID,
		hookID,
	).Scan(
		&hook.ID,
		&hook.ProjectID,
		&hook.Label,
		&hook.IsEnabled,
		&hook.Entity,
		&hook.EntityID,
		&hook.HookEvent,
		&hook.Headers,
		&hook.HookURL,
	)

	return &hook, errors.New(errors.UnknownError, err)
}

// UpdateHook updates all fields of a WebHook by project and hook ID
func (d *Database) UpdateHook(projectID, hookID string, hook *models.WebHook) *errors.DatastoreError {
	_, err := d.db.Exec(
		fmt.Sprintf(
			"UPDATE %s SET label=$1, isenabled=$2, entity=$3, entity_id=$4, hook_event=$5, headers=$6, hook_url=$7 WHERE id=$8 and project_id=$9",
			tableProjectWebHooks,
		),
		hook.Label,
		hook.IsEnabled,
		hook.Entity,
		hook.EntityID,
		hook.HookEvent,
		hook.Headers,
		hook.HookURL,
		hook.ID,
		projectID,
	)

	return errors.New(errors.UnknownError, err)
}

// DeleteHook permanently removes a WebHook by project and hook ID
func (d *Database) DeleteHook(projectID, hookID string) *errors.DatastoreError {
	_, err := d.db.Exec(
		fmt.Sprintf(
			"DELETE FROM %s WHERE id=$1 and project_id=$2",
			tableProjectWebHooks,
		),
		hookID,
		projectID,
	)
	return errors.New(errors.UnknownError, err)
}
