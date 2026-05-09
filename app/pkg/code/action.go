//go:generate codegen -type=int

package code

// Action: action service errors.
// Code must start with 1011xx.
const (
	// ErrRecordNotFound - 404: Record not found.
	ErrRecordNotFound int = iota + 101001

	// ErrMessageQuery - 500: Failed to query Message from Database.
	ErrMessageQuery

	// ErrMessageCreate - 500: Failed to create Message in Database.
	ErrMessageCreate
)
