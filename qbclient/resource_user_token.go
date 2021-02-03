package qbclient

import (
	"io"
	"net/http"
)

// CloneUserTokenInput models the input sent to POST /v1/usertoken/clone.
// See https://api.quickbase.com/v1/usertoken/clone
type CloneUserTokenInput struct {
	c *Client
	u string

	Name        string `json:"name" cliutil:"option=name"`
	Description string `json:"description" cliutil:"option=description"`
}

func (i *CloneUserTokenInput) url() string                  { return i.u }
func (i *CloneUserTokenInput) method() string               { return http.MethodPost }
func (i *CloneUserTokenInput) addHeaders(req *http.Request) { addHeadersJSON(req, i.c) }
func (i *CloneUserTokenInput) encode() ([]byte, error)      { return marshalJSON(i) }

// CloneUserTokenOutput models the output returned by POST /v1/usertoken/clone.
// See https://api.quickbase.com/v1/usertoken/clone
type CloneUserTokenOutput struct {
	ErrorProperties

	Active bool                       `json:"active"`
	Apps   []*CloneUserTokenOutputApp `json:"apps"`
	ID     int                        `json:"id"`
	Name   string                     `json:"name"`
	Token  string                     `json:"token"`

	// Description is defined in ErrorProperties.
}

func (o *CloneUserTokenOutput) decode(body io.ReadCloser) error { return unmarshalJSON(body, &o) }

// CloneUserTokenOutputApp models the apps property.
type CloneUserTokenOutputApp struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

// CloneUserToken sends a request to POST /v1/usertoken/clone.
// See https://api.quickbase.com/v1/usertoken/clone
func (c *Client) CloneUserToken(input *CloneUserTokenInput) (output *CloneUserTokenOutput, err error) {
	input.c = c
	input.u = c.URL + "/usertoken/clone"
	output = &CloneUserTokenOutput{}
	err = c.Do(input, output)
	return
}

// DeactivateUserTokenInput models the input sent to POST /v1/usertoken/deactivate.
// See https://developer.quickbase.com/operation/deactivateUserToken
type DeactivateUserTokenInput struct {
	c *Client
	u string

	Token string `json:"-" validate:"required" cliutil:"option=token"`
}

func (i *DeactivateUserTokenInput) url() string                  { return i.u }
func (i *DeactivateUserTokenInput) method() string               { return http.MethodPost }
func (i *DeactivateUserTokenInput) addHeaders(req *http.Request) { addHeadersJSON(req, i.c) }
func (i *DeactivateUserTokenInput) encode() ([]byte, error)      { return marshalJSON(i) }

// DeactivateUserTokenOutput models the output returned by POST /v1/usertoken/deactivate.
// See https://developer.quickbase.com/operation/deactivateUserToken
type DeactivateUserTokenOutput struct {
	ErrorProperties

	ID int `json:"id"`
}

func (o *DeactivateUserTokenOutput) decode(body io.ReadCloser) error { return unmarshalJSON(body, &o) }

// DeactivateUserToken sends a request to POST /v1/usertoken/deactivate.
// See https://developer.quickbase.com/operation/deactivateUserToken
func (c *Client) DeactivateUserToken(input *DeactivateUserTokenInput) (output *DeactivateUserTokenOutput, err error) {
	curToken := c.UserToken
	c.UserToken = input.Token
	defer func() { c.UserToken = curToken }()

	input.c = c
	input.u = c.URL + "/usertoken/deactivate"
	output = &DeactivateUserTokenOutput{}
	err = c.Do(input, output)
	return
}

// DeleteUserTokenInput models the input sent to DELETE /v1/usertoken.
// See https://developer.quickbase.com/operation/deleteUserToken
type DeleteUserTokenInput struct {
	c *Client
	u string

	Token string `json:"-" validate:"required" cliutil:"option=token"`
}

func (i *DeleteUserTokenInput) url() string                  { return i.u }
func (i *DeleteUserTokenInput) method() string               { return http.MethodDelete }
func (i *DeleteUserTokenInput) addHeaders(req *http.Request) { addHeadersJSON(req, i.c) }
func (i *DeleteUserTokenInput) encode() ([]byte, error)      { return marshalJSON(i) }

// DeleteUserTokenOutput models the output returned by DELETE /v1/usertoken.
// See https://developer.quickbase.com/operation/deleteUserToken
type DeleteUserTokenOutput struct {
	ErrorProperties

	ID int `json:"id"`
}

func (o *DeleteUserTokenOutput) decode(body io.ReadCloser) error { return unmarshalJSON(body, &o) }

// DeleteUserToken sends a request to DELETE /v1/usertoken.
// See https://developer.quickbase.com/operation/deleteUserToken
func (c *Client) DeleteUserToken(input *DeleteUserTokenInput) (output *DeleteUserTokenOutput, err error) {
	curToken := c.UserToken
	c.UserToken = input.Token
	defer func() { c.UserToken = curToken }()

	input.c = c
	input.u = c.URL + "/usertoken"
	output = &DeleteUserTokenOutput{}
	err = c.Do(input, output)
	return
}
