package postgres

import (
	"fmt"

	"github.com/anothrnick/machinable/dsi/models"
)

const tableAppTiers = "app_tiers"

// ListTiers retrieves all app tiers
func (d *Database) ListTiers() ([]*models.Tier, error) {
	rows, err := d.db.Query(
		fmt.Sprintf(
			"SELECT id, name, cost, requests, projects, storage FROM %s ORDER BY projects ASC",
			tableAppTiers,
		),
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	tiers := make([]*models.Tier, 0)
	for rows.Next() {
		tier := models.Tier{}
		err = rows.Scan(
			&tier.ID,
			&tier.Name,
			&tier.Cost,
			&tier.Requests,
			&tier.Projects,
			&tier.Storage,
		)
		if err != nil {
			return nil, err
		}

		tiers = append(tiers, &tier)
	}

	return tiers, rows.Err()
}
