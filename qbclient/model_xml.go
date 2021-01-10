package qbclient

import (
	"encoding/xml"
	"io"
	"net/http"
	"strings"

	"github.com/QuickBase/quickbase-cli/qberrors"
)

// XMLInput is implemented by requests to the XML API.
type XMLInput interface {
	Input

	// setAppToken sets the app token credential.
	setAppToken(token string)

	// setTicket sets the ticket credential.
	setTicket(token string)

	// setUserToken sets the user token credential.
	setUserToken(token string)
}

// XMLRequestParameters models common XML API request parameters.
type XMLRequestParameters struct {
	XMLName  xml.Name `xml:"qdbapi"`
	UserData string   `xml:"udata,omitempty"`
}

// XMLCredentialParameters models XML API credentials parameters and implements
// XMLInput.
type XMLCredentialParameters struct {
	AppToken  string `xml:"apptoken,omitempty"`
	Ticket    string `xml:"ticket,omitempty"`
	UserToken string `xml:"usertoken,omitempty"`
}

func (p *XMLCredentialParameters) setAppToken(token string)  { p.AppToken = token }
func (p *XMLCredentialParameters) setTicket(token string)    { p.Ticket = token }
func (p *XMLCredentialParameters) setUserToken(token string) { p.UserToken = token }

// XMLResponseParameters models common XML API response parameters and
// implements Output.
type XMLResponseParameters struct {
	XMLName     xml.Name `xml:"qdbapi" json:"-"`
	Action      string   `xml:"action" json:"-"`
	ErrorCode   int      `xml:"errcode" json:"errorCode,omitempty"`
	ErrorText   string   `xml:"errtext" json:"errorMessage,omitempty"`
	ErrorDetail string   `xml:"errdetail" json:"errorDetail,omitempty"`
	UserData    string   `xml:"udata,omitempty" json:"userData,omitempty"`
}

func (p *XMLResponseParameters) errorMessage() string { return p.ErrorText }
func (p *XMLResponseParameters) errorDetail() string  { return p.ErrorDetail }

// handleError handles XML API errors.
// See https://help.quickbase.com/api-guide/errorcodes.html
// See https://developer.mozilla.org/en-US/docs/Web/HTTP/Status
func (p *XMLResponseParameters) handleError(output Output, resp *http.Response) (err error) {
	if p.ErrorCode == 0 {
		p.ErrorText = ""
		return
	}

	if resp.StatusCode >= 400 && resp.StatusCode < 500 {
		msg := strings.ToLower(http.StatusText(resp.StatusCode))
		serr := qberrors.ErrSafe{Message: msg, StatusCode: resp.StatusCode}
		return qberrors.Client(serr)
	}

	if resp.StatusCode >= 500 {
		msg := strings.ToLower(http.StatusText(resp.StatusCode))
		serr := qberrors.ErrSafe{Message: msg, StatusCode: resp.StatusCode}
		return qberrors.Service(serr)
	}

	serr := qberrors.ErrSafe{Message: output.errorMessage()}
	switch p.ErrorCode {
	case 1:
		serr.StatusCode = http.StatusInternalServerError
	case 2:
		serr.StatusCode = p.statusFromCode2(output)
	case 3, 7, 19, 23, 24, 29, 34, 38, 70, 71, 73, 74, 77, 78, 114, 150, 151, 152:
		serr.StatusCode = http.StatusForbidden
	case 4, 13, 20, 21, 22, 27, 28, 83:
		serr.StatusCode = http.StatusUnauthorized
	case 5:
		serr.StatusCode = http.StatusNotAcceptable
	case 6, 8, 9, 10, 11, 12, 14, 15, 25, 26, 50, 51, 52, 53, 75, 76, 80, 87, 102, 103, 110:
		serr.StatusCode = http.StatusBadRequest
	case 30, 31, 32, 33, 35, 37, 54, 81, 112:
		serr.StatusCode = http.StatusNotFound
	case 36, 84, 85:
		serr.StatusCode = http.StatusInternalServerError
	case 60, 61:
		serr.StatusCode = http.StatusConflict
	case 82:
		serr.StatusCode = http.StatusGatewayTimeout
	case 100, 101, 105:
		serr.StatusCode = http.StatusServiceUnavailable
	case 104:
		serr.StatusCode = http.StatusTooManyRequests
	case 111, 113:
		serr.StatusCode = http.StatusUnprocessableEntity
	}

	err = qberrors.Client(nil).Safef(serr, "%s", output.errorDetail())
	return
}

// statusFromCode2 does its best to find a staus code for the error that
// resulted in a Quickbase error code 2.
func (p *XMLResponseParameters) statusFromCode2(output Output) int {
	if strings.Contains(p.ErrorDetail, "not found") {
		return http.StatusNotFound
	}

	return http.StatusBadRequest
}

// marshalXML marshals the API request into XML. This function is intended to
// be used in Input.marshal implementations.
func marshalXML(input XMLInput, c *Client) ([]byte, error) {
	// Add credentials.
	// TODO Support Tickets and App Tokens?
	// See https://github.com/QuickBase/quickbase-sdk-go/blob/master/creds.go#L26
	if c.UserToken != "" {
		input.setUserToken(c.UserToken)
	}

	return xml.Marshal(input)
}

func unmarshalXML(body io.ReadCloser, output Output) error {
	return xml.NewDecoder(body).Decode(&output)
}

// addHeadersJSON adds heads required for JSON requests.
func addHeadersXML(req *http.Request, c *Client, action string) {
	req.Header.Set("Content-Type", "application/xml")
	req.Header.Set("QUICKBASE-ACTION", action)
}
