package postgres

import (
	"fmt"
	"strings"
	"time"

	"github.com/anothrnick/machinable/dsi/models"
)

const tableProjectLogs = "project_logs"

// AddProjectLog saves a new log for a project
func (d *Database) AddProjectLog(projectID string, log *models.Log) error {
	_, err := d.db.Exec(
		fmt.Sprintf(
			"INSERT INTO %s (project_id, endpoint_type, verb, path, status_code, created, aligned, response_time, initiator, initiator_type, initiator_id, target_id) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)",
			tableProjectLogs,
		),
		projectID,
		log.EndpointType,
		log.Verb,
		log.Path,
		log.StatusCode,
		time.Unix(log.Created, 0),
		time.Unix(log.AlignedCreated, 0),
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
	args := make([]interface{}, 0)
	index := 1

	// query builders
	filterString := make([]string, 0)
	sortString := make([]string, 0)
	pageString := ""

	// projectID
	args = append(args, projectID)
	filterString = append(filterString, fmt.Sprintf("project_id=$%d", index))
	index++

	// valid filter/sort
	validFields := map[string]bool{"created": true, "initiator_type": true, "status_code": true, "endpoint_type": true}

	// filters
	filterErr := d.filterToQuery(filter, validFields, &filterString, &args, &index)
	if filterErr != nil {
		return nil, filterErr
	}

	// sort
	for key, val := range sort {
		// validate fields
		if _, ok := validFields[key]; !ok {
			// not a valid field, move on
			continue
		}
		direction := "DESC"
		if val > 0 {
			direction = "ASC"
		}
		sortString = append(sortString, fmt.Sprintf("%s %s", key, direction))
	}

	// paginate
	if limit >= 0 {
		args = append(args, limit)
		pageString += fmt.Sprintf(" LIMIT $%d", index)
		index++
	}

	if offset >= 0 {
		args = append(args, offset)
		pageString += fmt.Sprintf(" OFFSET $%d", index)
		index++
	}

	queryFields := "id, project_id, endpoint_type, verb, path, status_code, created, aligned, response_time, initiator, initiator_type, initiator_id, target_id"
	orderBy := ""
	if len(sortString) > 0 {
		orderBy = fmt.Sprintf(" ORDER BY %s", strings.Join(sortString, ", "))
	}

	query := fmt.Sprintf(
		"SELECT %s FROM %s WHERE %s%s%s",
		queryFields,
		tableProjectLogs,
		strings.Join(filterString, " AND "),
		orderBy,
		pageString,
	)

	rows, err := d.db.Query(
		query,
		args...,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	logs := make([]*models.Log, 0)
	for rows.Next() {
		log := models.Log{}
		created := time.Time{}
		aligned := time.Time{}
		err = rows.Scan(
			&log.ID,
			&log.ProjectID,
			&log.EndpointType,
			&log.Verb,
			&log.Path,
			&log.StatusCode,
			&created,
			&aligned,
			&log.ResponseTime,
			&log.Initiator,
			&log.InitiatorType,
			&log.InitiatorID,
			&log.TargetID,
		)
		log.Created = created.Unix()
		log.AlignedCreated = aligned.Unix()
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
	args := make([]interface{}, 0)
	index := 1

	// query builders
	filterString := make([]string, 0)

	// projectID
	args = append(args, projectID)
	filterString = append(filterString, fmt.Sprintf("project_id=$%d", index))
	index++

	// valid filter/sort
	validFields := map[string]bool{"created": true, "initiator_type": true, "status_code": true}

	// filters
	filterErr := d.filterToQuery(filter, validFields, &filterString, &args, &index)
	if filterErr != nil {
		return 0, filterErr
	}

	query := fmt.Sprintf(
		"SELECT count(id) FROM %s WHERE %s",
		tableProjectLogs,
		strings.Join(filterString, " AND "),
	)

	err := d.db.QueryRow(
		query,
		args...,
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
