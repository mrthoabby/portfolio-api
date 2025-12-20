package types

// ErrNotFound is returned when a record is not found in the store.
type ErrNotFound struct {
	Message string
}

func (e ErrNotFound) Error() string {
	if e.Message == "" {
		return "record not found"
	}
	return e.Message
}

// IsNotFoundError checks if the error is a not found error.
func IsNotFoundError(err error) bool {
	_, ok := err.(ErrNotFound)
	return ok
}

