package interfaces

import "github.com/machinable/machinable/dsi/models"

// TiersDatastore exposes functions for app tiers
type TiersDatastore interface {
	ListTiers() ([]*models.Tier, error)
}
