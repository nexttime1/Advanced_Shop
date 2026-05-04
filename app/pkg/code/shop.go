//go:generate codegen -type=int

package code

// Auth: authentication and authorization errors.
// Code must start with 1003xx.
const (
	// ErrUnauthorized - 401: User not logged in.
	ErrUnauthorized int = iota + 100301

	// ErrInvalidUserID - 400: Invalid user ID format.
	ErrInvalidUserID

	// ErrRoleNotConfigured - 500: User role not configured.
	ErrRoleNotConfigured

	// ErrInvalidRoleFormat - 400: Invalid role format.
	ErrInvalidRoleFormat

	// ErrInvalidRole - 403: Invalid role value.
	ErrInvalidRole
)
