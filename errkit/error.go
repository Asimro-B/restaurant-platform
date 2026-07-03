package errkit

import "errors"

func New(message string) error {
	return errors.New("api.restaurant: " + message)
}

var (
	ErrInvalidRequestBody       = errors.New("invalid request body")
	ErrInvalidData        error = errors.New("api.restaurant:Invalid Data")
	ErrUserNotFound       error = errors.New(
		"api.restaurant:User Not Found. Please check your email",
	)
	ErrInvalidCredentials = errors.New("api.restaurant:Invalid credential")
)
