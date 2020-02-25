package postgres

import (
	"fmt"

	"github.com/machinable/machinable/dsi/models"
)

const tableAppUsers = "app_users"

// GetAppUserByUsername attempts to find a user by username, if the user does not exist the user will be nil, error will be !nil
func (d *Database) GetAppUserByUsername(userName string) (*models.User, error) {
	user := &models.User{}

	err := d.db.QueryRow(
		fmt.Sprintf(
			"SELECT id, email, username, password_hash, created, active from %s WHERE username=$1",
			tableAppUsers,
		),
		userName,
	).Scan(
		&user.ID,
		&user.Email,
		&user.Username,
		&user.PasswordHash,
		&user.Created,
		&user.Active,
	)

	return user, err
}

// GetAppUserByID attempts to find a user by ID, if the user does not exist the user will be nil, error will be !nil
func (d *Database) GetAppUserByID(userID string) (*models.User, error) {
	user := &models.User{}

	err := d.db.QueryRow(
		fmt.Sprintf(
			"SELECT id, email, username, password_hash, created, tier_id, active from %s WHERE id=$1",
			tableAppUsers,
		),
		userID,
	).Scan(
		&user.ID,
		&user.Email,
		&user.Username,
		&user.PasswordHash,
		&user.Created,
		&user.Tier,
		&user.Active,
	)

	return user, err
}

// CreateAppUser saves a new application user, updates user ID
func (d *Database) CreateAppUser(user *models.User) error {
	err := d.db.QueryRow(
		fmt.Sprintf(
			"INSERT INTO %s (email, username, password_hash, created) VALUES ($1, $2, $3, $4) RETURNING id",
			tableAppUsers,
		),
		user.Email,
		user.Username,
		user.PasswordHash,
		user.Created,
	).Scan(&user.ID)

	return err
}

// UpdateUserPassword updates the user's password
func (d *Database) UpdateUserPassword(userID, passwordHash string) error {
	_, err := d.db.Exec(
		fmt.Sprintf(
			"UPDATE %s SET password_hash=$1 WHERE id=$2",
			tableAppUsers,
		),
		passwordHash,
		userID,
	)

	return err
}

// ActivateUser updates the user active field
func (d *Database) ActivateUser(userID string, active bool) error {
	_, err := d.db.Exec(
		fmt.Sprintf(
			"UPDATE %s SET active=$1 WHERE id=$2",
			tableAppUsers,
		),
		active,
		userID,
	)

	return err
}
