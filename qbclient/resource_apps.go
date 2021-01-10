package qbclient

import (
	"encoding/json"
	"io"
	"net/http"
	"net/url"
)

// CreateAppInput models the input sent to POST /v1/apps.
// See https://developer.quickbase.com/operation/createApp
type CreateAppInput struct {
	c *Client
	u string

	Name            string      `json:"name" validate:"required" cliutil:"option=name usage='name of the app'"`
	Description     string      `json:"description,omitempty" cliutil:"option=description usage='description of the app'"`
	AssignUserToken bool        `json:"assignToken,omitempty" cliutil:"option=assign-token usage='assign the user token to the app'"`
	Variable        []*Variable `json:"variables,omitempty"`
}

func (i *CreateAppInput) url() string                  { return i.u }
func (i *CreateAppInput) method() string               { return http.MethodPost }
func (i *CreateAppInput) addHeaders(req *http.Request) { addHeadersJSON(req, i.c) }
func (i *CreateAppInput) encode() ([]byte, error)      { return marshalJSON(i) }

// CreateAppOutput models the output returned by POST /v1/apps.
// See https://developer.quickbase.com/operation/createApp
type CreateAppOutput struct {
	ErrorProperties
	App
}

func (o *CreateAppOutput) decode(body io.ReadCloser) error { return unmarshalJSON(body, &o) }

// CreateApp sends a request to POST /v1/apps.
// See https://developer.quickbase.com/operation/getApp
func (c *Client) CreateApp(input *CreateAppInput) (output *CreateAppOutput, err error) {
	input.c = c
	input.u = c.URL + "/apps"
	output = &CreateAppOutput{}
	err = c.Do(input, output)
	return
}

// GetAppInput models the input sent to GET /v1/apps/{appId}.
// See https://developer.quickbase.com/operation/getApp
type GetAppInput struct {
	c *Client
	u string

	AppID string `json:"-" validate:"required" cliutil:"option=app-id"`
}

func (i *GetAppInput) url() string                  { return i.u }
func (i *GetAppInput) method() string               { return http.MethodGet }
func (i *GetAppInput) addHeaders(req *http.Request) { addHeadersJSON(req, i.c) }
func (i *GetAppInput) encode() ([]byte, error)      { return marshalJSON(i) }

// GetAppOutput models the output returned by GET /v1/apps/{appId}.
// See https://developer.quickbase.com/operation/getApp
type GetAppOutput struct {
	ErrorProperties
	App
}

func (o *GetAppOutput) decode(body io.ReadCloser) error { return unmarshalJSON(body, &o) }

// GetApp sends a request to GET /v1/apps/{appId}.
// See https://developer.quickbase.com/operation/getApp
func (c *Client) GetApp(input *GetAppInput) (output *GetAppOutput, err error) {
	input.c = c
	input.u = c.URL + "/apps/" + url.PathEscape(input.AppID)
	output = &GetAppOutput{}
	err = c.Do(input, output)
	return
}

// GetAppByID sends a request to GET /v1/apps/{appId} and gets an app by ID.
// See https://developer.quickbase.com/operation/getApp
func (c *Client) GetAppByID(id string) (*GetAppOutput, error) {
	return c.GetApp(&GetAppInput{AppID: id})
}

// UpdateAppInput models the input sent to POST /v1/apps/{appId}.
// See https://developer.quickbase.com/operation/updateApp
type UpdateAppInput struct {
	c *Client
	u string

	AppID       string      `json:"-" validate:"required" cliutil:"option=app-id"`
	Name        string      `json:"name,omitempty" cliutil:"option=name usage='name of the app'"`
	Description string      `json:"description,omitempty" cliutil:"option=description usage='description of the app'"`
	Variable    []*Variable `json:"variables,omitempty"`
}

func (i *UpdateAppInput) url() string                  { return i.u }
func (i *UpdateAppInput) method() string               { return http.MethodPost }
func (i *UpdateAppInput) addHeaders(req *http.Request) { addHeadersJSON(req, i.c) }
func (i *UpdateAppInput) encode() ([]byte, error)      { return marshalJSON(i) }

// UpdateAppOutput models the output returned by POST /v1/apps/{appId}.
// See https://developer.quickbase.com/operation/updateApp
type UpdateAppOutput struct {
	ErrorProperties
	App
}

func (o *UpdateAppOutput) decode(body io.ReadCloser) error { return unmarshalJSON(body, &o) }

// UpdateApp sends a request to POST /v1/apps/{appId}.
// See https://developer.quickbase.com/operation/updateApp
func (c *Client) UpdateApp(input *UpdateAppInput) (output *UpdateAppOutput, err error) {
	input.c = c
	input.u = c.URL + "/apps/" + url.PathEscape(input.AppID)
	output = &UpdateAppOutput{}
	err = c.Do(input, output)
	return
}

// DeleteAppInput models the input sent to DELETE /v1/apps/{appId}.
// See https://developer.quickbase.com/operation/deleteApp
type DeleteAppInput struct {
	c *Client
	u string

	AppID string `json:"-" validate:"required" cliutil:"option=app-id"`
	Name  string `json:"name" validate:"required" cliutil:"option=name usage='name of the app'"`
}

func (i *DeleteAppInput) url() string                  { return i.u }
func (i *DeleteAppInput) method() string               { return http.MethodDelete }
func (i *DeleteAppInput) addHeaders(req *http.Request) { addHeadersJSON(req, i.c) }
func (i *DeleteAppInput) encode() ([]byte, error)      { return marshalJSON(i) }

// DeleteAppOutput models the output returned by DELETE /v1/apps/{appId}.
// See https://developer.quickbase.com/operation/deleteApp
type DeleteAppOutput struct {
	ErrorProperties

	ID string `json:"deletedAppId,omitempty"`
}

func (o *DeleteAppOutput) decode(body io.ReadCloser) error { return unmarshalJSON(body, &o) }

// DeleteApp sends a request to DELETE /v1/apps/{appId}.
// See https://developer.quickbase.com/operation/deleteApp
func (c *Client) DeleteApp(input *DeleteAppInput) (output *DeleteAppOutput, err error) {
	input.c = c
	input.u = c.URL + "/apps/" + url.PathEscape(input.AppID)
	output = &DeleteAppOutput{}
	err = c.Do(input, output)
	return
}

// ListAppEventsInput models the input sent to GET /v1/apps/{appId}/events.
// See https://developer.quickbase.com/operation/getAppEvents
type ListAppEventsInput struct {
	c *Client
	u string

	AppID string `json:"-" validate:"required" cliutil:"option=app-id"`
}

func (i *ListAppEventsInput) url() string                  { return i.u }
func (i *ListAppEventsInput) method() string               { return http.MethodGet }
func (i *ListAppEventsInput) addHeaders(req *http.Request) { addHeadersJSON(req, i.c) }
func (i *ListAppEventsInput) encode() ([]byte, error)      { return marshalJSON(i) }

// ListAppEventsOutput models the output returned by GET /v1/apps/{appId}/events.
// See https://developer.quickbase.com/operation/getAppEvents
type ListAppEventsOutput struct {
	ErrorProperties

	Events []*ListAppEventsOutputEvent `json:"events,omitempty"`
}

func (o *ListAppEventsOutput) decode(body io.ReadCloser) error { return unmarshalJSON(body, &o) }

// UnmarshalJSON implements json.UnmarshalJSON by unmarshaling the payload into
// ListTablesOutput.Events.
func (o *ListAppEventsOutput) UnmarshalJSON(b []byte) (err error) {
	var v []*ListAppEventsOutputEvent
	if err = json.Unmarshal(b, &v); err == nil {
		o.Events = v
	} else {
		err = json.Unmarshal(b, &o.ErrorProperties)
	}
	return
}

// MarshalJSON implements json.MarshalJSON by marshaling output.Tables.
func (o *ListAppEventsOutput) MarshalJSON() ([]byte, error) {
	return json.Marshal(o.Events)
}

// ListAppEventsOutputEvent models the event object.
type ListAppEventsOutputEvent struct {
	Type     string `json:"type"`
	Owner    *User  `json:"owner"`
	IsActive bool   `json:"isActive"`
	TableID  string `json:"tableId"`
	Name     string `json:"name"`
	URL      string `json:"url,omitempty"`
}

// ListAppEvents sends a request to GET /v1/apps/{appId}/events.
// See https://developer.quickbase.com/operation/getAppEvents
func (c *Client) ListAppEvents(input *ListAppEventsInput) (output *ListAppEventsOutput, err error) {
	input.c = c
	input.u = c.URL + "/apps/" + url.PathEscape(input.AppID) + "/events"
	output = &ListAppEventsOutput{}
	err = c.Do(input, output)
	return
}

// CopyAppInput models the input sent to POST /v1/apps/{appId}/copy.
// See https://developer.quickbase.com/operation/copyApp
type CopyAppInput struct {
	c *Client
	u string

	AppID       string                  `json:"-" validate:"required" cliutil:"option=app-id usage='unique identifier of an app'"`
	Name        string                  `json:"name" validate:"required" cliutil:"option=name usage='name of the app'"`
	Description string                  `json:"description,omitempty" cliutil:"option=description usage='description of the app'"`
	Properties  *CopyAppInputProperties `json:"properties,omitempty"`
}

func (i *CopyAppInput) url() string                  { return i.u }
func (i *CopyAppInput) method() string               { return http.MethodPost }
func (i *CopyAppInput) addHeaders(req *http.Request) { addHeadersJSON(req, i.c) }
func (i *CopyAppInput) encode() ([]byte, error)      { return marshalJSON(i) }

// CopyAppInputProperties models the properties property.
type CopyAppInputProperties struct {
	AssignUserToken   bool `json:"assignUserToken,omitempty" cliutil:"option=assign-token usage='assign the user token to the app'"`
	ExcludeFiles      bool `json:"excludeFiles,omitempty" cliutil:"option=exclude-files usage='exclude attached files if --copy-data is passed'"`
	KeepData          bool `json:"keepData,omitempty" cliutil:"option=keep-data usage='copy data'"`
	KeepUsersAndRoles bool `json:"usersAndRoles,omitempty" cliutil:"option=keep-users-roles usage='copy users and roles'"`
}

// CopyAppOutput models the output returned by POST /v1/apps/{appId}/copy.
// See https://developer.quickbase.com/operation/copyApp
type CopyAppOutput struct {
	ErrorProperties
	App

	AncestorID string `json:"ancestorId,omitempty"`
}

func (o *CopyAppOutput) decode(body io.ReadCloser) error { return unmarshalJSON(body, &o) }

// CopyApp sends a request to POST /v1/apps/{appId}/copy.
// See https://developer.quickbase.com/operation/copyApp
func (c *Client) CopyApp(input *CopyAppInput) (output *CopyAppOutput, err error) {
	input.c = c
	input.u = c.URL + "/apps/" + url.PathEscape(input.AppID) + "/copy"
	output = &CopyAppOutput{}
	err = c.Do(input, output)
	return
}
