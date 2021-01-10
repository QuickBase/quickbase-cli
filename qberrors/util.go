package qberrors

import (
	"errors"
	"strings"
)

// IsSafe returns true if err is safe to show the user.
func IsSafe(err error) bool { return errors.As(err, &ErrSafe{}) }

// SafeMessage returns an error message that is safe for the user to see.
func SafeMessage(err error) string {
	if IsSafe(err) {
		if serr, ok := err.(ErrSafe); ok {
			return serr.Message
		}
		return SafeMessage(errors.Unwrap(err))
	}
	return defaultMsg
}

// SafeDetail returns detail about the error that is safe for the user to see.
func SafeDetail(err error) string {
	if IsSafe(err) {
		s := err.Error()
		m := SafeMessage(err)
		return strings.TrimRight(s[0:len(s)-len(m)], ": ")
	}
	return ""
}

// StatusCode returns the status code associated with the error.
func StatusCode(err error) int {
	if serr, ok := err.(ErrSafe); ok {
		if serr.StatusCode != 0 {
			return serr.StatusCode
		}
	} else if err != nil {
		return StatusCode(errors.Unwrap(err))
	}
	return defaultCode
}

// Upstream returns the error chain that caused the user to be shown an error.
// The error chain is assumed to be unsafe for the user to see.
func Upstream(err error) error {
	if IsSafe(err) {
		return upstream(err)
	}
	return err
}

func upstream(err error) error {
	if ierr, ok := err.(Error); ok {
		return ierr.Upstream()
	} else if werr := errors.Unwrap(err); werr != nil {
		return upstream(werr)
	}
	return nil
}

// NotFoundError returns an ErrClient with ErrClient.safe set as NotFound and
// additional context according to the format specifier.
func NotFoundError(format string, a ...interface{}) error {
	return Client(nil).Safef(NotFound, format, a...)
}

// InvalidInputError returns an ErrClient with ErrClient.safe set as BadRequest
// and additional context according to the format specifier.
func InvalidInputError(err error, format string, a ...interface{}) error {
	return Client(err).Safef(BadRequest, format, a...)
}
