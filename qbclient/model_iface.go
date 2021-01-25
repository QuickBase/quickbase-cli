package qbclient

import (
	"io"
	"net/http"
)

// Input models the payload of API requests.
type Input interface {

	// url returns the URL the API request is sent to.
	url() string

	// method is the HTTP method used when sending API requests.
	method() string

	// addHeaders adds HTTP headers to the API request.
	addHeaders(req *http.Request)

	// encode encodes the request and writes it to io.Writer.
	encode() ([]byte, error)
}

// Output models the payload of API responses.
type Output interface {

	// decode parses the response in io.ReadCloser to Output.
	decode(io.ReadCloser) error

	// errorMessage returns the error message, if any.
	errorMessage() string

	// errorDetail returns the error detail, if any.
	errorDetail() string

	// handleError handles errors returned by the API.
	handleError(Output, *http.Response) error
}
