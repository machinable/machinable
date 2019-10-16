package interfaces

import "github.com/anothrnick/machinable/dsi/models"

// TiersDatastore exposes functions for app tiers
type TiersDatastore interface {
	ListTiers() ([]*models.Tier, error)
}
