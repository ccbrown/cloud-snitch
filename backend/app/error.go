package app

import (
	"errors"
	"net/http"

	"github.com/stripe/stripe-go/v81"
)

type UserFacingError interface {
	error
	UserFacingError() string
}

type userErrorMessage string

func (e userErrorMessage) UserFacingError() string {
	return string(e)
}

func (e userErrorMessage) Error() string {
	return string(e)
}

func NewUserError(message string) UserFacingError {
	return userErrorMessage(message)
}

type InternalError struct {
	Cause error
}

func (e InternalError) UserFacingError() string {
	return e.Error()
}

func (e InternalError) Error() string {
	return "Internal error."
}

func (sess *Session) SanitizedError(err error) UserFacingError {
	if err == nil {
		return nil
	} else if ufe, ok := err.(UserFacingError); ok {
		return ufe
	}
	sess.logger.Error(err.Error())
	return InternalError{Cause: err}
}

type AuthorizationError struct{}

func (e AuthorizationError) UserFacingError() string {
	return e.Error()
}

func (e AuthorizationError) Error() string {
	return "Unauthorized."
}

type AuthenticationError struct{}

func (e AuthenticationError) UserFacingError() string {
	return e.Error()
}

func (e AuthenticationError) Error() string {
	return "Bad authorization."
}

type NotFoundError string

func (e NotFoundError) UserFacingError() string {
	return e.Error()
}

func (e NotFoundError) Error() string {
	return string(e)
}

func IsStripeBadRequestError(err error) bool {
	if err == nil {
		return false
	}
	var stripeErr *stripe.Error
	if errors.As(err, &stripeErr) {
		return stripeErr.HTTPStatusCode == http.StatusBadRequest
	}
	return false
}
