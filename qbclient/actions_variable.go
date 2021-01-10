package qbclient

import (
	"io"
	"net/http"
	"net/url"
)

// GetVariableInput models the XML API request sent to API_GetDBvar
// See https://help.quickbase.com/api-guide/index.html#getdbvar.html
type GetVariableInput struct {
	XMLRequestParameters
	XMLCredentialParameters

	c *Client
	u string

	AppID string `xml:"-" validate:"required" cliutil:"option=app-id"`
	Name  string `xml:"varname" validate:"required" cliutil:"option=name"`
}

func (i *GetVariableInput) method() string               { return http.MethodPost }
func (i *GetVariableInput) url() string                  { return i.u }
func (i *GetVariableInput) addHeaders(req *http.Request) { addHeadersXML(req, i.c, "API_GetDBvar") }
func (i *GetVariableInput) encode() ([]byte, error)      { return marshalXML(i, i.c) }

// GetVariableOutput models the XML API response returned by API_GetDBvar.
// See https://help.quickbase.com/api-guide/index.html#getdbvar.html
type GetVariableOutput struct {
	XMLResponseParameters
	Variable
}

func (o *GetVariableOutput) decode(body io.ReadCloser) error { return unmarshalXML(body, o) }

// GetVariable sends an XML API request to API_GetDBvar.
// See https://help.quickbase.com/api-guide/index.html#getdbvar.html
func (c *Client) GetVariable(input *GetVariableInput) (output *GetVariableOutput, err error) {
	input.c = c
	input.u = "https://" + url.PathEscape(c.ReamlHostname) + "/db/" + url.PathEscape(input.AppID)

	output = &GetVariableOutput{}
	err = c.Do(input, output)
	if err == nil {
		output.Name = input.Name
	}

	return
}

// SetVariableInput models a request sent to API_SetDBvar via the XML API.
// See https://help.quickbase.com/api-guide/index.html#setdbvar.html
type SetVariableInput struct {
	XMLRequestParameters
	XMLCredentialParameters

	c *Client
	u string

	AppID string `xml:"-" validate:"required" cliutil:"option=app-id"`
	Name  string `xml:"varname" validate:"required" cliutil:"option=name"`
	Value string `xml:"value" cliutil:"option=value"`
}

func (i *SetVariableInput) method() string               { return http.MethodPost }
func (i *SetVariableInput) url() string                  { return i.u }
func (i *SetVariableInput) addHeaders(req *http.Request) { addHeadersXML(req, i.c, "API_SetDBvar") }
func (i *SetVariableInput) encode() ([]byte, error)      { return marshalXML(i, i.c) }

// SetVariableOutput models the XML API response returned by API_SetDBvar
// See https://help.quickbase.com/api-guide/index.html#setdbvar.html
type SetVariableOutput struct {
	XMLResponseParameters

	Name  string `xml:"-" json:"name,omitempty"`
	Value string `xml:"-" json:"value,omitempty"`
}

func (o *SetVariableOutput) decode(body io.ReadCloser) error { return unmarshalXML(body, o) }

// SetVariable sends an XML API request to API_SetDBvar.
// See https://help.quickbase.com/api-guide/index.html#setdbvar.html
func (c *Client) SetVariable(input *SetVariableInput) (output *SetVariableOutput, err error) {
	input.c = c
	input.u = "https://" + url.PathEscape(c.ReamlHostname) + "/db/" + url.PathEscape(input.AppID)

	output = &SetVariableOutput{}
	err = c.Do(input, output)
	if err == nil {
		output.Name = input.Name
		output.Value = input.Value
	}

	return
}
