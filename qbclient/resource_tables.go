package qbclient

import (
	"encoding/json"
	"io"
	"net/http"
	"net/url"
)

// CreateTableInput models the input sent to POST /v1/tables?appId={appId}.
// See https://developer.quickbase.com/operation/createTable
type CreateTableInput struct {
	c *Client
	u string

	AppID        string `json:"-" validate:"required" cliutil:"option=app-id"`
	Name         string `json:"name" validate:"required" cliutil:"option=name"`
	Description  string `json:"description" cliutil:"option=description"`
	IconName     string `json:"iconName" cliutil:"option=icon-name"`
	SingularNoun string `json:"singularNoun" cliutil:"option=singular-noun"`
	PluralNoun   string `json:"pluralNoun" cliutil:"option=plural-noun"`
}

func (i *CreateTableInput) url() string                  { return i.u }
func (i *CreateTableInput) method() string               { return http.MethodPost }
func (i *CreateTableInput) addHeaders(req *http.Request) { addHeadersJSON(req, i.c) }
func (i *CreateTableInput) encode() ([]byte, error)      { return marshalJSON(i) }

// CreateTableOutput models the output returned by POST /v1/tables?appId={appId}.
// See https://developer.quickbase.com/operation/createTable
type CreateTableOutput struct {
	ErrorProperties

	TableID      string `json:"id,omitempty"`
	Name         string `json:"name,omitempty"`
	IconName     string `json:"iconName,omitempty"`
	SingularNoun string `json:"singularNoun,omitempty"`
	PluralNoun   string `json:"pluralNoun,omitempty"`

	// Description is defined in ErrorProperties.
}

func (o *CreateTableOutput) decode(body io.ReadCloser) error { return unmarshalJSON(body, &o) }

// CreateTable sends a request to POST /v1/tables?appId={appId}.
// See https://developer.quickbase.com/operation/createTable
func (c *Client) CreateTable(input *CreateTableInput) (output *CreateTableOutput, err error) {
	input.c = c
	input.u = c.URL + "/tables?appId=" + url.QueryEscape(input.AppID)
	output = &CreateTableOutput{}
	err = c.Do(input, output)
	return
}

// ListTablesInput models the input sent to GET /v1/tables?appId={appId}.
// See https://developer.quickbase.com/operation/getAppTables
type ListTablesInput struct {
	c *Client
	u string

	AppID string `json:"-" validate:"required" cliutil:"option=app-id"`
}

func (i *ListTablesInput) url() string                  { return i.u }
func (i *ListTablesInput) method() string               { return http.MethodGet }
func (i *ListTablesInput) addHeaders(req *http.Request) { addHeadersJSON(req, i.c) }
func (i *ListTablesInput) encode() ([]byte, error)      { return marshalJSON(i) }

// ListTablesOutput models the input sent to GET /v1/tables?appId={appId}.
// See https://developer.quickbase.com/operation/getAppTables
type ListTablesOutput struct {
	ErrorProperties

	Tables []*ListTablesOutputTable `json:"tables,omitempty"`
}

func (o *ListTablesOutput) decode(body io.ReadCloser) error { return unmarshalJSON(body, &o) }

// UnmarshalJSON implements json.UnmarshalJSON by unmarshaling the payload into
// ListTablesOutput.Tables.
func (o *ListTablesOutput) UnmarshalJSON(b []byte) (err error) {
	var v []*ListTablesOutputTable
	if err = json.Unmarshal(b, &v); err == nil {
		o.Tables = v
	} else {
		err = json.Unmarshal(b, &o.ErrorProperties)
	}
	return
}

// MarshalJSON implements json.MarshalJSON by marshaling output.Tables.
func (o *ListTablesOutput) MarshalJSON() ([]byte, error) {
	return json.Marshal(o.Tables)
}

// ListTablesOutputTable models the table object.
type ListTablesOutputTable struct {
	Name               string     `json:"name"`
	TableID            string     `json:"id"`
	Alias              string     `json:"alias"`
	Description        string     `json:"description"`
	Created            *Timestamp `json:"created"`
	Updated            *Timestamp `json:"updated"`
	NextRecordID       int        `json:"nextRecordId"`
	NextFieldID        int        `json:"nextFieldId"`
	DefaultSortFieldID int        `json:"defaultSortFieldId"`
	DefaultSortOrder   string     `json:"defaultSortOrder"`
	KeyFieldID         int        `json:"keyFieldId"`
	SingleRecordName   string     `json:"singleRecordName"`
	PluralRecordName   string     `json:"pluralRecordName"`
	TimeZone           string     `json:"timeZone"`
	DateFormat         string     `json:"dateFormat"`
	SizeLimit          string     `json:"sizeLimit"`
	SpaceRemaining     string     `json:"spaceRemaining"`
	SpaceUsed          string     `json:"spaceUsed"`
}

// ListTables sends a request to GET /v1/tables?appId={appId}.
// See https://developer.quickbase.com/operation/getAppTables
func (c *Client) ListTables(input *ListTablesInput) (output *ListTablesOutput, err error) {
	input.c = c
	input.u = c.URL + "/tables?appId=" + url.QueryEscape(input.AppID)
	output = &ListTablesOutput{}
	err = c.Do(input, output)
	return
}

// ListTablesByAppID sends a request to GET /v1/tables?appId={appId} and gets a
// list of tables in an app by its ID.
// See https://developer.quickbase.com/operation/getAppTables
func (c *Client) ListTablesByAppID(id string) (*ListTablesOutput, error) {
	return c.ListTables(&ListTablesInput{AppID: id})
}

// GetTableInput models the input sent to GET /v1/tables/{tableId}?appId={appId}s.
// See https://developer.quickbase.com/operation/getTable
type GetTableInput struct {
	c *Client
	u string

	AppID   string `json:"-" validate:"required" cliutil:"option=app-id"`
	TableID string `json:"-" validate:"required" cliutil:"option=table-id"`
}

func (i *GetTableInput) url() string                  { return i.u }
func (i *GetTableInput) method() string               { return http.MethodGet }
func (i *GetTableInput) addHeaders(req *http.Request) { addHeadersJSON(req, i.c) }
func (i *GetTableInput) encode() ([]byte, error)      { return marshalJSON(i) }

// GetTableOutput models the output returned by GET /v1/tables/{tableId}?appId={appId}.
// See https://developer.quickbase.com/operation/getTable
type GetTableOutput struct {
	ErrorProperties

	Name               string     `json:"name,omitempty"`
	TableID            string     `json:"id,omitempty"`
	Alias              string     `json:"alias,omitempty"`
	Created            *Timestamp `json:"created,omitempty"`
	Updated            *Timestamp `json:"updated,omitempty"`
	NextRecordID       int        `json:"nextRecordId,omitempty"`
	NextFieldID        int        `json:"nextFieldId,omitempty"`
	DefaultSortFieldID int        `json:"defaultSortFieldId,omitempty"`
	DefaultSortOrder   string     `json:"defaultSortOrder,omitempty"`
	KeyFieldID         int        `json:"keyFieldId,omitempty"`
	SingleRecordName   string     `json:"singleRecordName,omitempty"`
	PluralRecordName   string     `json:"pluralRecordName,omitempty"`
	TimeZone           string     `json:"timeZone,omitempty"`
	DateFormat         string     `json:"dateFormat,omitempty"`

	// Description is defined in ErrorProperties.
}

func (o *GetTableOutput) decode(body io.ReadCloser) error { return unmarshalJSON(body, &o) }

// GetTable sends a request to GET /v1/tables/{tableId}?appId={appId}.
// See https://developer.quickbase.com/operation/getTable
func (c *Client) GetTable(input *GetTableInput) (output *GetTableOutput, err error) {
	input.c = c
	input.u = c.URL + "/tables/" + url.PathEscape(input.TableID) + "?appId=" + url.QueryEscape(input.AppID)
	output = &GetTableOutput{}
	err = c.Do(input, output)
	return
}

// UpdateTableInput models the input sent to POST /v1/tables/{tableId}?appId={appId}.
// See https://developer.quickbase.com/operation/updateTable
type UpdateTableInput struct {
	c *Client
	u string

	AppID        string `json:"-" validate:"required" cliutil:"option=app-id"`
	TableID      string `json:"-" validate:"required" cliutil:"option=table-id"`
	Name         string `json:"name,omitempty" cliutil:"option=name"`
	Description  string `json:"description,omitempty" cliutil:"option=description"`
	IconName     string `json:"iconName,omitempty" cliutil:"option=icon-name"`
	SingularNoun string `json:"singularNoun,omitempty" cliutil:"option=singular-noun"`
	PluralNoun   string `json:"pluralNoun,omitempty" cliutil:"option=plural-noun"`
}

func (i *UpdateTableInput) url() string                  { return i.u }
func (i *UpdateTableInput) method() string               { return http.MethodPost }
func (i *UpdateTableInput) addHeaders(req *http.Request) { addHeadersJSON(req, i.c) }
func (i *UpdateTableInput) encode() ([]byte, error)      { return marshalJSON(i) }

// UpdateTableOutput models the output returned by POST /v1/tables/{tableId}?appId={appId}.
// See https://developer.quickbase.com/operation/updateTable
type UpdateTableOutput struct {
	ErrorProperties

	TableID      string `json:"id,omitempty"`
	Name         string `json:"name,omitempty"`
	IconName     string `json:"iconName,omitempty"`
	SingularNoun string `json:"singularNoun,omitempty"`
	PluralNoun   string `json:"pluralNoun,omitempty"`

	// Description is defined in ErrorProperties.
}

func (o *UpdateTableOutput) decode(body io.ReadCloser) error { return unmarshalJSON(body, &o) }

// UpdateTable sends a request to POST /v1/tables/{tableId}?appId={appId}.
// See https://developer.quickbase.com/operation/updateTable
func (c *Client) UpdateTable(input *UpdateTableInput) (output *UpdateTableOutput, err error) {
	input.c = c
	input.u = c.URL + "/tables/" + url.PathEscape(input.TableID) + "?appId=" + url.QueryEscape(input.AppID)
	output = &UpdateTableOutput{}
	err = c.Do(input, output)
	return
}

// DeleteTableInput models the input sent to DELETE /v1/tables/{tableId}?appId={appId}.
// See https://developer.quickbase.com/operation/deleteTable
type DeleteTableInput struct {
	c *Client
	u string

	AppID   string `json:"-" validate:"required" cliutil:"option=app-id"`
	TableID string `json:"-" validate:"required" cliutil:"option=table-id"`
}

func (i *DeleteTableInput) url() string                  { return i.u }
func (i *DeleteTableInput) method() string               { return http.MethodDelete }
func (i *DeleteTableInput) addHeaders(req *http.Request) { addHeadersJSON(req, i.c) }
func (i *DeleteTableInput) encode() ([]byte, error)      { return marshalJSON(i) }

// DeleteTableOutput models the output returned by DELETE /v1/tables/{tableId}?appId={appId}.
// See https://developer.quickbase.com/operation/deleteTable
type DeleteTableOutput struct {
	ErrorProperties

	TableID string `json:"deletedTableId,omitempty"`
}

func (o *DeleteTableOutput) decode(body io.ReadCloser) error { return unmarshalJSON(body, &o) }

// DeleteTable sends a request to DELETE /v1/tables/{tableId}?appId={appId}.
// See https://developer.quickbase.com/operation/deleteTable
func (c *Client) DeleteTable(input *DeleteTableInput) (output *DeleteTableOutput, err error) {
	input.c = c
	input.u = c.URL + "/tables/" + url.PathEscape(input.TableID) + "?appId=" + url.QueryEscape(input.AppID)
	output = &DeleteTableOutput{}
	err = c.Do(input, output)
	return
}
