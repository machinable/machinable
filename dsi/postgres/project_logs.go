package postgres

import "github.com/anothrnick/machinable/dsi/models"

const tableProjectLogs = "project_logs"

// AddProjectLog saves a new log for a project
func (d *Database) AddProjectLog(project string, log *models.Log) error {
	return nil
}

// ListProjectLogs retrieves logs based on the limit, offset, filter, and sort parameters
func (d *Database) ListProjectLogs(project string, limit, offset int64, filter *models.Filters, sort map[string]int) ([]*models.Log, error) {
	return nil, nil
}

// CountProjectLogs returns the count of logs for a project
func (d *Database) CountProjectLogs(project string, filter *models.Filters) (int64, error) {
	return 0, nil
}

// DropProjectLogs drops the collection for this project's logs
func (d *Database) DropProjectLogs(project string) error {
	return nil
}
