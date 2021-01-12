package qbclient

// TODO Replace this file with https://github.com/go-playground/validator

import (
	"errors"
	"fmt"
	"regexp"
)

var (
	reHostname *regexp.Regexp
)

func init() {
	reHostname = regexp.MustCompile(`^(([a-zA-Z0-9]|[a-zA-Z0-9][a-zA-Z0-9\-]*[a-zA-Z0-9])\.)+([A-Za-z0-9]|[A-Za-z0-9][A-Za-z0-9\-]*[A-Za-z0-9])$`)
}

// ValidateStringFn validates a string.
type ValidateStringFn func(string) error

// NoValidation always returns a nil error.
func NoValidation(s string) error {
	return nil
}

// ValidateNotEmptyFn returns a function that validates a string isn't empty.
func ValidateNotEmptyFn(label string) ValidateStringFn {
	return func(s string) (err error) {
		if s == "" {
			err = fmt.Errorf("%s required", label)
		}
		return
	}
}

// ValidateHostname validates the passed hostname.
func ValidateHostname(hostname string) error {
	if hostname == "" {
		return errors.New("hostname required")
	} else if !reHostname.MatchString(hostname) {
		return errors.New("invalid hostname, expecting format example.quickbase.com")
	}
	return nil
}
