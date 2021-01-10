package qbclient

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/QuickBase/quickbase-cli/qberrors"
)

// ErrorProperties contains properties returned during errors.
type ErrorProperties struct {
	Message     string `json:"message,omitempty"`
	Description string `json:"description,omitempty"`
}

func (p *ErrorProperties) errorMessage() string { return p.Message }
func (p *ErrorProperties) errorDetail() string  { return p.Description }

func (p *ErrorProperties) handleError(output Output, resp *http.Response) (err error) {
	switch true {
	case resp.StatusCode >= 200 && resp.StatusCode < 300:
		return
	case resp.StatusCode >= 400 && resp.StatusCode < 500:
		serr := qberrors.ErrSafe{Message: output.errorMessage(), StatusCode: resp.StatusCode}
		return qberrors.Client(serr).Safef(serr, "%s", output.errorDetail())
	default:
		serr := qberrors.ErrSafe{Message: output.errorMessage(), StatusCode: resp.StatusCode}
		return qberrors.Service(serr).Safef(serr, "%s", output.errorDetail())
	}
}

// addHeadersJSON adds heads required for JSON requests.
func addHeadersJSON(req *http.Request, c *Client) {
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("QB-Realm-Hostname", c.ReamlHostname)
	req.Header.Add("User-Agent", c.UserAgent)

	var authstr string
	if c.TemporaryToken != "" {
		authstr = fmt.Sprintf("QB-TEMP-TOKEN %s", c.TemporaryToken)
	} else if c.UserToken != "" {
		authstr = fmt.Sprintf("QB-USER-TOKEN %s", c.UserToken)
	}
	if authstr != "" {
		req.Header.Add("Authorization", authstr)
	}
}

// marshalJSON marshals the API request into JSON. This function is intended to
// be used in Input.marshal implementations.
func marshalJSON(input interface{}) (b []byte, err error) {
	if b, err = json.Marshal(input); err == nil {
		if bytes.Equal(b, []byte(`{}`)) {
			b = []byte(``)
		}
	}
	return
}

// unmarshalJSON unmarshals a JSON API response into an Output. This function
// is intended to be used in Output.unmarshal implementations.
func unmarshalJSON(body io.ReadCloser, output interface{}) error {
	return json.NewDecoder(body).Decode(&output)
}
