//go:generate codegen -type=int

package code

// User: user service errors.
// Code must start with 1004xx.
const (
	// ErrUserNotFound - 404: User not found.
	ErrUserNotFound int = iota + 100401

	// ErrUserPasswordIncorrect - 401: User password incorrect.
	ErrUserPasswordIncorrect

	// ErrCodeNotExist - 404: Verification code not exist.
	ErrCodeNotExist

	// ErrCodeInCorrect - 400: Verification code incorrect.
	ErrCodeInCorrect

	// ErrUserAlreadyExists - 400: User already exists.
	ErrUserAlreadyExists

	// ErrSmsSend - 500: Failed to send SMS.
	ErrSmsSend

	// ErrForbidden - 403: User privilege insufficient.
	ErrForbidden
)
