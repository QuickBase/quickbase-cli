package qbclient

import (
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"strconv"
)

// ListFieldsInput models the input sent to GET /v1/fields?tableId={tableId}.
// See https://developer.quickbase.com/operation/getFields
type ListFieldsInput struct {
	c *Client
	u string

	TableID                 string `json:"-" validate:"required" cliutil:"option=table-id"`
	IncludeFieldPermissions bool   `json:"includeFieldPerms" cliutil:"option=include-field-permissions"`
}

func (i *ListFieldsInput) url() string                  { return i.u }
func (i *ListFieldsInput) method() string               { return http.MethodGet }
func (i *ListFieldsInput) addHeaders(req *http.Request) { addHeadersJSON(req, i.c) }
func (i *ListFieldsInput) encode() ([]byte, error)      { return marshalJSON(i) }

// ListFieldsOutput models the output returned by GET /v1/fields?tableId={tableId}.
// See https://developer.quickbase.com/operation/getFields
type ListFieldsOutput struct {
	ErrorProperties

	Fields []*ListFieldsOutputField `json:"fields,omitempty"`
}

func (o *ListFieldsOutput) decode(body io.ReadCloser) error { return unmarshalJSON(body, &o) }

// UnmarshalJSON implements json.UnmarshalJSON by unmarshaling the payload into
// ListFieldsOutput.Fields.
func (o *ListFieldsOutput) UnmarshalJSON(b []byte) (err error) {
	var v []*ListFieldsOutputField
	if err = json.Unmarshal(b, &o.ErrorProperties); err == nil {
		return
	}
	if err = json.Unmarshal(b, &v); err == nil {
		o.Fields = v
	}
	return
}

// ListFieldsOutputField models the field object.
type ListFieldsOutputField struct {
	Field
	FieldID    int                              `json:"id,omitempty"`
	Properties *ListFieldsOutputFieldProperties `json:"properties,omitempty"`
}

// ListFieldsOutputFieldProperties models the field object properties.
type ListFieldsOutputFieldProperties struct {
	FieldProperties
}

// ListFields sends a request to GET /v1/fields?tableId={tableId}.
// See https://developer.quickbase.com/operation/getFields
func (c *Client) ListFields(input *ListFieldsInput) (output *ListFieldsOutput, err error) {
	input.c = c
	input.u = c.URL + "/fields?tableId=" + url.QueryEscape(input.TableID)
	output = &ListFieldsOutput{}
	err = c.Do(input, output)
	return
}

// ListFieldsByTableID sends a request to GET /v1/fields?tableId={tableId} an
// lists fields for the passed table.
// See https://developer.quickbase.com/operation/getFields
func (c *Client) ListFieldsByTableID(tableID string) (*ListFieldsOutput, error) {
	return c.ListFields(&ListFieldsInput{TableID: tableID})
}

// CreateFieldInput models the input sent to POST /v1/fields?tableId={tableId}.
// See https://developer.quickbase.com/operation/createField
type CreateFieldInput struct {
	Field

	c *Client
	u string

	TableID    string                      `json:"-" validate:"required" cliutil:"option=table-id"`
	Properties *CreateFieldInputProperties `json:"properties,omitempty"`
}

func (i *CreateFieldInput) url() string                  { return i.u }
func (i *CreateFieldInput) method() string               { return http.MethodPost }
func (i *CreateFieldInput) addHeaders(req *http.Request) { addHeadersJSON(req, i.c) }
func (i *CreateFieldInput) encode() ([]byte, error)      { return marshalJSON(i) }

// CreateFieldInputProperties models the "properties" property.
type CreateFieldInputProperties struct {
	FieldProperties
}

// CreateFieldOutput models the output returned by POST /v1/fields?tableId={tableId}.
// See https://developer.quickbase.com/operation/createField
type CreateFieldOutput struct {
	ErrorProperties
	Field

	FieldID    int                          `json:"id"`
	Properties *CreateFieldOutputProperties `json:"properties,omitempty"`
}

func (o *CreateFieldOutput) decode(body io.ReadCloser) error { return unmarshalJSON(body, &o) }

// CreateFieldOutputProperties models the "properties" property.
type CreateFieldOutputProperties struct {
	FieldProperties
}

// CreateField sends a request to POST /v1/fields?tableId={tableId}.
// See https://developer.quickbase.com/operation/createField
func (c *Client) CreateField(input *CreateFieldInput) (output *CreateFieldOutput, err error) {
	input.c = c
	input.u = c.URL + "/fields?tableId=" + url.QueryEscape(input.TableID)
	input.Create = true
	output = &CreateFieldOutput{}
	err = c.Do(input, output)
	return
}

// DeleteFieldsInput models the input sent to DELETE v1/fields?tableId={tableId}.
// See https://developer.quickbase.com/operation/deleteFields
type DeleteFieldsInput struct {
	c *Client
	u string

	TableID  string `json:"-" validate:"required" cliutil:"option=table-id"`
	FieldIDs []int  `json:"fieldIds" validate:"required,min=1" cliutil:"option=field-id"`
}

func (i *DeleteFieldsInput) url() string                  { return i.u }
func (i *DeleteFieldsInput) method() string               { return http.MethodDelete }
func (i *DeleteFieldsInput) addHeaders(req *http.Request) { addHeadersJSON(req, i.c) }
func (i *DeleteFieldsInput) encode() ([]byte, error)      { return marshalJSON(i) }

// DeleteFieldsOutput models the output returned by DELETE v1/fields?tableId={tableId}.
// See https://developer.quickbase.com/operation/deleteFields
type DeleteFieldsOutput struct {
	ErrorProperties

	DeletedFieldIDs []int    `json:"deletedFieldIds,omitempty"`
	Errors          []string `json:"errors,omitempty"`
}

func (o *DeleteFieldsOutput) decode(body io.ReadCloser) error { return unmarshalJSON(body, &o) }

// DeleteFields sends a request to DELETE v1/fields?tableId={tableId}.
// See https://developer.quickbase.com/operation/deleteFields
func (c *Client) DeleteFields(input *DeleteFieldsInput) (output *DeleteFieldsOutput, err error) {
	input.c = c
	input.u = c.URL + "/fields?tableId=" + url.QueryEscape(input.TableID)
	output = &DeleteFieldsOutput{}
	err = c.Do(input, output)
	return
}

// GetFieldInput models the input sent to GET /v1/fields/{fieldId}?tableId={tableId}.
// See https://developer.quickbase.com/operation/getField
type GetFieldInput struct {
	c *Client
	u string

	TableID string `json:"-" validate:"required" cliutil:"option=table-id"`
	FieldID int    `json:"-" validate:"required" cliutil:"option=field-id"`
}

func (i *GetFieldInput) url() string                  { return i.u }
func (i *GetFieldInput) method() string               { return http.MethodGet }
func (i *GetFieldInput) addHeaders(req *http.Request) { addHeadersJSON(req, i.c) }
func (i *GetFieldInput) encode() ([]byte, error)      { return marshalJSON(i) }

// GetFieldOutput models the output returned by GET /v1/fields/{fieldId}?tableId={tableId}.
// See https://developer.quickbase.com/operation/getField
type GetFieldOutput struct {
	ErrorProperties
	Field

	FieldID    int                       `json:"id,omitempty"`
	Properties *GetFieldOutputProperties `json:"properties,omitempty"`
}

func (o *GetFieldOutput) decode(body io.ReadCloser) error { return unmarshalJSON(body, &o) }

// GetFieldOutputProperties models the "properties" property.
type GetFieldOutputProperties struct {
	FieldProperties
}

// GetField sends a request to GET /v1/fields/{fieldId}?tableId={tableId}.
// See https://developer.quickbase.com/operation/getField
func (c *Client) GetField(input *GetFieldInput) (output *GetFieldOutput, err error) {
	input.c = c
	fieldID := strconv.Itoa(input.FieldID)
	input.u = c.URL + "/fields/" + url.PathEscape(fieldID) + "?tableId=" + url.QueryEscape(input.TableID)
	output = &GetFieldOutput{}
	err = c.Do(input, output)
	return
}

// GetFieldByID sends a request to GET /v1/fields/{fieldId}?tableId={tableId}
// and gets a field by its ID.
// See https://developer.quickbase.com/operation/getField
func (c *Client) GetFieldByID(fid int) (*GetFieldOutput, error) {
	return c.GetField(&GetFieldInput{FieldID: fid})
}

// UpdateFieldInput models the input sent to POST /v1/fields/{fieldId}?tableId={tableId}.
// See https://developer.quickbase.com/operation/updateField
type UpdateFieldInput struct {
	Field

	c *Client
	u string

	TableID    string                      `json:"-" validate:"required" cliutil:"option=table-id"`
	FieldID    int                         `json:"-" validate:"required" cliutil:"option=field-id"`
	Properties *UpdateFieldInputProperties `json:"properties,omitempty"`
}

func (i *UpdateFieldInput) url() string                  { return i.u }
func (i *UpdateFieldInput) method() string               { return http.MethodPost }
func (i *UpdateFieldInput) addHeaders(req *http.Request) { addHeadersJSON(req, i.c) }
func (i *UpdateFieldInput) encode() ([]byte, error)      { return marshalJSON(i) }

// UpdateFieldInputProperties models the "properties" property.
type UpdateFieldInputProperties struct {
	FieldProperties
}

// UpdateFieldOutput models the output returned by POST /v1/fields/{fieldId}?tableId={tableId}.
// See https://developer.quickbase.com/operation/updateField
type UpdateFieldOutput struct {
	ErrorProperties
	Field

	FieldID    int                          `json:"id"`
	Properties *CreateFieldOutputProperties `json:"properties,omitempty"`
}

func (o *UpdateFieldOutput) decode(body io.ReadCloser) error { return unmarshalJSON(body, &o) }

// UpdateField sends a request to POST /v1/fields/{fieldId}?tableId={tableId}.
// See https://developer.quickbase.com/operation/updateField
func (c *Client) UpdateField(input *UpdateFieldInput) (output *UpdateFieldOutput, err error) {
	input.c = c
	fieldID := strconv.Itoa(input.FieldID)
	input.u = c.URL + "/fields/" + url.PathEscape(fieldID) + "?tableId=" + url.QueryEscape(input.TableID)
	input.Create = false
	output = &UpdateFieldOutput{}
	err = c.Do(input, output)
	return
}
