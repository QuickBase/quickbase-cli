package qbclient

import "net/http"

// Plugin is implemented by plugins that intercept the HTTP request and
// response when consuming the Quick Base API.
type Plugin interface {
	PreRequest(req *http.Request)
	PostResponse(resp *http.Response)
}
