package postgres

import (
	"fmt"
	"time"

	"github.com/anothrnick/machinable/dsi/models"
)

const tableProjectLogs = "project_logs"

// AddProjectLog saves a new log for a project
func (d *Database) AddProjectLog(projectID string, log *models.Log) error {
	_, err := d.db.Exec(
		fmt.Sprintf(
			"INSERT INTO %s (project_id, endpoint_type, verb, path, status_code, created, response_time, initiator, initiator_type, initiator_id, target_id) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)",
			tableProjectLogs,
		),
		projectID,
		log.EndpointType,
		log.Verb,
		log.Path,
		log.StatusCode,
		time.Now(),
		log.ResponseTime,
		log.Initiator,
		log.InitiatorType,
		log.InitiatorID,
		log.TargetID,
	)

	return err
}

// ListProjectLogs retrieves logs based on the limit, offset, filter, and sort parameters
func (d *Database) ListProjectLogs(projectID string, limit, offset int64, filter *models.Filters, sort map[string]int) ([]*models.Log, error) {
	rows, err := d.db.Query(
		fmt.Sprintf(
			"SELECT id, project_id, endpoint_type, verb, path, status_code, created, response_time, initiator, initiator_type, initiator_id, target_id FROM %s WHERE project_id=$1",
			tableProjectLogs,
		),
		projectID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	logs := make([]*models.Log, 0)
	for rows.Next() {
		log := models.Log{}
		created := time.Time{}
		err = rows.Scan(
			&log.ID,
			&log.ProjectID,
			&log.EndpointType,
			&log.Verb,
			&log.Path,
			&log.StatusCode,
			&created,
			&log.ResponseTime,
			&log.Initiator,
			&log.InitiatorType,
			&log.InitiatorID,
			&log.TargetID,
		)
		log.Created = created.Unix()
		if err != nil {
			return nil, err
		}

		logs = append(logs, &log)
	}

	return logs, rows.Err()
}

// CountProjectLogs returns the count of logs for a project
func (d *Database) CountProjectLogs(projectID string, filter *models.Filters) (int64, error) {
	var count int64
	err := d.db.QueryRow(
		fmt.Sprintf(
			"SELECT count(id) FROM %s WHERE project_id=$1",
			tableProjectLogs,
		),
		projectID,
	).Scan(&count)
	return count, err
}

// DropProjectLogs drops the collection for this project's logs
func (d *Database) DropProjectLogs(projectID string) error {
	_, err := d.db.Exec(
		fmt.Sprintf(
			"DELETE FROM %s WHERE project_id=$1",
			tableProjectLogs,
		),
		projectID,
	)
	return err
}
