package code

const (
	// ErrUserNotFound - 404: User not found.
	ErrUserNotFound int = iota + 100401

	// ErrUserPasswordIncorrect - 404: User Password Incorrect.
	ErrUserPasswordIncorrect

	// ErrCodeNotExist - 404: Code Not Exist.
	ErrCodeNotExist

	// ErrCodeInCorrect - 404: Code Not Correct.
	ErrCodeInCorrect
)
