package qbclient

import (
	"io"
	"net/http"
)

// InsertRecordsInput models the input sent to POST /v1/records.
// See https://developer.quickbase.com/operation/upsert
type InsertRecordsInput struct {
	c *Client
	u string

	Data           []map[int]*InsertRecordsInputData `json:"data" validate:"required,min=1" cliutil:"option=data func=record"`
	To             string                            `json:"to" validate:"required" cliutil:"option=to"`
	MergeFieldID   int                               `json:"mergeFieldId,omitempty" cliutil:"option=merge-field-id"`
	FieldsToReturn []int                             `json:"fieldsToReturn,omitempty" cliutil:"option=fields-to-return"`
}

func (i *InsertRecordsInput) url() string                  { return i.u }
func (i *InsertRecordsInput) method() string               { return http.MethodPost }
func (i *InsertRecordsInput) addHeaders(req *http.Request) { addHeadersJSON(req, i.c) }
func (i *InsertRecordsInput) encode() ([]byte, error)      { return marshalJSON(i) }

// InsertRecordsInputData models the data property.
type InsertRecordsInputData struct {
	Value *Value `json:"value" validate:"required"`
}

// InsertRecordsOutput models the output returned by POST /v1/records.
// See https://developer.quickbase.com/operation/upsert
type InsertRecordsOutput struct {
	ErrorProperties

	Metadata *InsertRecordsOutputMetadata `json:"metadata,omitempty"`
}

func (o *InsertRecordsOutput) decode(body io.ReadCloser) error { return unmarshalJSON(body, &o) }

// InsertRecordsOutputMetadata models the metadata property.
type InsertRecordsOutputMetadata struct {
	CreatedRecordIDs              []int `json:"createdRecordIds"`
	TotalNumberOfRecordsProcessed int   `json:"totalNumberOfRecordsProcessed"`
	UnchangedRecordIDs            []int `json:"unchangedRecordIds"`
	IpdatedRecordIDs              []int `json:"updatedRecordIds"`
}

// InsertRecords sends a request to POST /v1/records.
// See https://developer.quickbase.com/operation/upsert
func (c *Client) InsertRecords(input *InsertRecordsInput) (output *InsertRecordsOutput, err error) {
	input.c = c
	input.u = c.URL + "/records"
	output = &InsertRecordsOutput{}
	err = c.Do(input, output)
	return
}

// DeleteRecordsInput models the input sent to DELETE /v1/records.
// See https://developer.quickbase.com/operation/deleteRecords
type DeleteRecordsInput struct {
	c *Client
	u string

	From  string `json:"from" validate:"required" cliutil:"option=from"`
	Where string `json:"where" validate:"required" cliutil:"option=where func=query"`
}

func (i *DeleteRecordsInput) url() string                  { return i.u }
func (i *DeleteRecordsInput) method() string               { return http.MethodDelete }
func (i *DeleteRecordsInput) addHeaders(req *http.Request) { addHeadersJSON(req, i.c) }
func (i *DeleteRecordsInput) encode() ([]byte, error)      { return marshalJSON(i) }

// DeleteRecordsOutput models the output returned by DELETE /v1/records.
// See https://developer.quickbase.com/operation/deleteRecords
type DeleteRecordsOutput struct {
	ErrorProperties

	NumberDeleted int `json:"numberDeleted,omitempty"`
}

func (o *DeleteRecordsOutput) decode(body io.ReadCloser) error { return unmarshalJSON(body, &o) }

// DeleteRecords sends a request to DELETE /v1/records.
// See https://developer.quickbase.com/operation/deleteRecords
func (c *Client) DeleteRecords(input *DeleteRecordsInput) (output *DeleteRecordsOutput, err error) {
	input.c = c
	input.u = c.URL + "/records"
	output = &DeleteRecordsOutput{}
	err = c.Do(input, output)
	return
}

// QueryRecordsInput models the input sent to POST /v1/records/query.
// See https://developer.quickbase.com/operation/runQuery
type QueryRecordsInput struct {
	c *Client
	u string

	Select  []int                       `json:"select" validate:"required,min=1" cliutil:"option=select"`
	From    string                      `json:"from" validate:"required" cliutil:"option=from"`
	Where   string                      `json:"where" cliutil:"option=where func=query"`
	GroupBy []*QueryRecordsInputGroupBy `json:"groupBy,omitempty" cliutil:"option=group-by func=group"`
	SortBy  []*QueryRecordsInputSortBy  `json:"sortBy,omitempty" cliutil:"option=sort-by func=sort"`
}

func (i *QueryRecordsInput) url() string                  { return i.u }
func (i *QueryRecordsInput) method() string               { return http.MethodPost }
func (i *QueryRecordsInput) addHeaders(req *http.Request) { addHeadersJSON(req, i.c) }
func (i *QueryRecordsInput) encode() ([]byte, error)      { return marshalJSON(i) }

// QueryRecordsInputGroupBy models the groupBy objects.
type QueryRecordsInputGroupBy struct {
	FieldID  int    `json:"fieldId"`
	Grouping string `json:"grouping"`
}

// QueryRecordsInputSortBy models the sortBy objects.
type QueryRecordsInputSortBy struct {
	FieldID int    `json:"fieldId"`
	Order   string `json:"order"`
}

// QueryRecordsOutput models the output returned by POST /v1/records/query.
// See https://developer.quickbase.com/operation/runQuery
type QueryRecordsOutput struct {
	ErrorProperties

	Data     []map[int]*QueryRecordsOutputData `json:"data,omitempty"`
	Fields   []*QueryRecordsOutputFields       `json:"fields,omitempty"`
	Metadata *QueryRecordsOutputMetadata       `json:"metadata,omitempty"`
}

func (o *QueryRecordsOutput) decode(body io.ReadCloser) error { return unmarshalJSON(body, &o) }

// QueryRecordsOutputData models objects in the data property.
type QueryRecordsOutputData struct {
	Value *Value `json:"value"`
}

// QueryRecordsOutputFields models the objects in the fields property.
type QueryRecordsOutputFields struct {
	FieldID int    `json:"id"`
	Label   string `json:"label"`
	Type    string `json:"type"`
}

// QueryRecordsOutputMetadata models the metadata property.
type QueryRecordsOutputMetadata struct {
	TotalRecords int `json:"totalRecords"`
	NumRecords   int `json:"numRecords"`
	NumFields    int `json:"numFields"`
	Skip         int `json:"skip"`
}

// QueryRecords sends a request to POST /v1/records/query.
// See https://developer.quickbase.com/operation/runQuery
func (c *Client) QueryRecords(input *QueryRecordsInput) (output *QueryRecordsOutput, err error) {
	input.c = c
	input.u = c.URL + "/records/query"
	output = &QueryRecordsOutput{}
	err = c.Do(input, output)
	return
}
