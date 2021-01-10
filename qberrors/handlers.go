package qberrors

import (
	"encoding/json"
	"net/http"
)

// HandleErrorJSON handles a JSON unmarshaling error for input passed by a user
// by normalizing messages and returning either a ErrClient or ErrInternal.
func HandleErrorJSON(err error) error {
	if serr, ok := err.(*json.SyntaxError); ok {
		return Client(serr).Safef(BadRequest, "%s: offset %v", serr, serr.Offset)
	}

	if terr, ok := err.(*json.UnmarshalTypeError); ok {
		return Client(terr).Safef(
			BadRequest,
			"%s field expects %s value, %s passed: offset %v",
			terr.Field,
			terr.Type,
			terr.Value,
			terr.Offset,
		)
	}

	serr := ErrSafe{"internal error decoding json", http.StatusInternalServerError}
	return Internal(err).Safe(serr)
}

// HandleErrorValidation handles github.com/go-playground/validator validation
// errors for input passed by a user and returns an ErrClient.
func HandleErrorValidation(err error) error {
	return Client(nil).Safef(BadRequest, "input not valid: %s", err)
}
