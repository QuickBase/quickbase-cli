// Package qberrors provides common errors to standardize error handling across
// Quickbae golang projects.
//
// See https://blog.golang.org/go1.13-errors
package qberrors

import (
	"fmt"
	"net/http"
)

const (
	defaultMsg  = "internal error"
	defaultCode = http.StatusInternalServerError
)

// ErrSafe errors.
var (
	BadRequest        = ErrSafe{"bad request", http.StatusBadRequest}
	InvalidDataType   = ErrSafe{"data type not valid", http.StatusBadRequest}
	InvalidInput      = ErrSafe{"input not valid", http.StatusBadRequest}
	InvalidSyntax     = ErrSafe{"syntax not valid", http.StatusBadRequest}
	NotFound          = ErrSafe{"not found", http.StatusNotFound}
	ServceUnavailable = ErrSafe{"service unavailable", http.StatusServiceUnavailable}
)

// ErrSafe is an error that is assumed to be safe to show to the user. Errors
// that wrap ErrSafe are also assumed to be safe to show the user, inclusive of
// all subsequent wraps up the chain.
type ErrSafe struct {
	Message    string
	StatusCode int
}

func (e ErrSafe) Error() string { return e.Message }

// SafeErrorf returns a wrapped ErrSafe given the format specifier.
func SafeErrorf(err error, format string, a ...interface{}) error {
	format += ": %w"
	a = append(a, err)
	return fmt.Errorf(format, a...)
}

// Error is implemented by errors.
type Error interface {

	// Upstream returns the error chain that caused showing the user an error.
	// It is assumed that the error chain is unsafe to show the user.
	Upstream() error

	// Retry returns whether the operation should be retried.
	Retry() bool
}

// ErrClient is an error due to client input.
// The operation should not be retried.
type ErrClient struct {
	upstream error
	safe     error
}

// Client returns an ErrClient.
func Client(err error) *ErrClient { return &ErrClient{err, BadRequest} }

func (e ErrClient) Unwrap() error { return e.safe }
func (e ErrClient) Error() string { return e.safe.Error() }

// Upstream implements Error.Upstream.
func (e ErrClient) Upstream() error { return e.upstream }

// Retry implements Error.Retry.
func (e ErrClient) Retry() bool { return false }

// Safe sets ErrClient.safe as err.
// TODO add variadic argument that wraps err
func (e *ErrClient) Safe(err error) error {
	e.safe = err
	return e
}

// Safef sets ErrClient.safe as a wrapped err.
func (e *ErrClient) Safef(err error, format string, a ...interface{}) error {
	e.safe = SafeErrorf(err, format, a...)
	return e
}

// ErrInternal is an error with the application.
// The operation should not be retried.
type ErrInternal struct {
	upstream error
	safe     error
}

// Internal returns an ErrInternal.
func Internal(err error) *ErrInternal {
	return &ErrInternal{err, ErrSafe{defaultMsg, defaultCode}}
}

func (e ErrInternal) Unwrap() error { return e.safe }
func (e ErrInternal) Error() string { return e.safe.Error() }

// Upstream implements Error.Upstream.
func (e ErrInternal) Upstream() error { return e.upstream }

// Retry implements Error.Retry.
func (e ErrInternal) Retry() bool { return false }

// Safe sets ErrInternal.safe as err.
// TODO add variadic argument that wraps err
func (e *ErrInternal) Safe(err error) error {
	e.safe = err
	return e
}

// Safef sets ErrInternal.safe as a wrapped err.
func (e *ErrInternal) Safef(err error, format string, a ...interface{}) error {
	e.safe = SafeErrorf(err, format, a...)
	return e
}

// ErrService is an error connecting to a dependent service.
// The operation can be retried.
type ErrService struct {
	upstream error
	safe     error
}

// Service returns an ErrService.
func Service(err error) *ErrService { return &ErrService{err, ServceUnavailable} }

func (e ErrService) Unwrap() error { return e.safe }
func (e ErrService) Error() string { return e.safe.Error() }

// Upstream implements Error.Upstream.
func (e ErrService) Upstream() error { return e.upstream }

// Retry implements Error.Retry.
func (e ErrService) Retry() bool { return true }

// Safe sets ErrService.safe as err.
// TODO add variadic argument that wraps err
func (e *ErrService) Safe(err error) error {
	e.safe = err
	return e
}

// Safef sets ErrService.safe as a wrapped err.
func (e *ErrService) Safef(err error, format string, a ...interface{}) error {
	e.safe = SafeErrorf(err, format, a...)
	return e
}
