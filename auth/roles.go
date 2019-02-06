package auth

const (
	// RoleUser is a constant value for the user role
	RoleUser = "user"
	// RoleAdmin is a constant value for the admin role
	RoleAdmin = "admin"
)

// ValidRoles is a list of allowed roles for project users and api keys
var ValidRoles = [2]string{RoleUser, RoleAdmin}
