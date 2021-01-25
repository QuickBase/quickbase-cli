package qbclient

import (
	"encoding/json"
	"io"
	"net/http"
	"net/url"
)

// GetReportInput models the input sent to GET /v1/reports/{reportId}?tableId={tableId}.
// See https://developer.quickbase.com/operation/getReport
type GetReportInput struct {
	c *Client
	u string

	TableID  string `json:"-" validate:"required" cliutil:"option=table-id"`
	ReportID string `json:"-" validate:"required" cliutil:"option=report-id"`
}

func (i *GetReportInput) url() string                  { return i.u }
func (i *GetReportInput) method() string               { return http.MethodGet }
func (i *GetReportInput) addHeaders(req *http.Request) { addHeadersJSON(req, i.c) }
func (i *GetReportInput) encode() ([]byte, error)      { return marshalJSON(i) }

// GetReportOutput models the output returned by GET /v1/reports/{reportId}?tableId={tableId}.
// See https://developer.quickbase.com/operation/getReport
type GetReportOutput struct {
	ErrorProperties
	Report
}

func (o *GetReportOutput) decode(body io.ReadCloser) error { return unmarshalJSON(body, &o) }

// GetReport sends a request to GET /v1/reports/{reportId}?tableId={tableId}.
// See https://developer.quickbase.com/operation/getReport
func (c *Client) GetReport(input *GetReportInput) (output *GetReportOutput, err error) {
	input.c = c
	input.u = c.URL + "/reports/" + url.PathEscape(input.ReportID) + "?tableId=" + url.QueryEscape(input.TableID)
	output = &GetReportOutput{}
	err = c.Do(input, output)
	return
}

// ListReportsInput models the input sent to GET /v1/reports?tableId={tableId}.
// See https://developer.quickbase.com/operation/GetReportReports
type ListReportsInput struct {
	c *Client
	u string

	TableID string `json:"-" validate:"required" cliutil:"option=table-id"`
}

func (i *ListReportsInput) url() string                  { return i.u }
func (i *ListReportsInput) method() string               { return http.MethodGet }
func (i *ListReportsInput) addHeaders(req *http.Request) { addHeadersJSON(req, i.c) }
func (i *ListReportsInput) encode() ([]byte, error)      { return marshalJSON(i) }

// ListReportsOutput models the input sent to GET /v1/reports?tableId={tableId}.
// See https://developer.quickbase.com/operation/GetReportReports
type ListReportsOutput struct {
	ErrorProperties

	Reports []*ReportWithDescripton
}

func (o *ListReportsOutput) decode(body io.ReadCloser) error { return unmarshalJSON(body, &o) }

// UnmarshalJSON implements json.UnmarshalJSON by unmarshaling the payload into
// ListTablesOutput.Tables.
func (o *ListReportsOutput) UnmarshalJSON(b []byte) (err error) {
	var v []*ReportWithDescripton
	if err = json.Unmarshal(b, &v); err == nil {
		o.Reports = v
	} else {
		err = json.Unmarshal(b, &o.ErrorProperties)
	}
	return
}

// MarshalJSON implements json.MarshalJSON by marshaling output.Tables.
func (o *ListReportsOutput) MarshalJSON() ([]byte, error) {
	return json.Marshal(o.Reports)
}

// ListReports sends a request to GET /v1/reports?tableId={tableId}.
// See https://developer.quickbase.com/operation/GetReportReports
func (c *Client) ListReports(input *ListReportsInput) (output *ListReportsOutput, err error) {
	input.c = c
	input.u = c.URL + "/reports?tableId=" + url.QueryEscape(input.TableID)
	output = &ListReportsOutput{}
	err = c.Do(input, output)
	return
}

// RunReportInput models the input sent to POST /v1/reports/{reportId}/run?tableId={tableId}.
// See https://developer.quickbase.com/operation/runReport
type RunReportInput struct {
	c *Client
	u string

	TableID  string `json:"-" validate:"required" cliutil:"option=table-id"`
	ReportID string `json:"-" validate:"required" cliutil:"option=report-id"`
	Skip     int    `json:"-" cliutil:"option=skip"`
	Top      int    `json:"-" cliutil:"option=top"`
}

func (i *RunReportInput) url() string                  { return i.u }
func (i *RunReportInput) method() string               { return http.MethodPost }
func (i *RunReportInput) addHeaders(req *http.Request) { addHeadersJSON(req, i.c) }
func (i *RunReportInput) encode() ([]byte, error)      { return marshalJSON(i) }

// RunReportOutput models the output returned by POST /v1/reports/{reportId}/run?tableId={tableId}.
// See https://developer.quickbase.com/operation/runReport
type RunReportOutput struct {
	ErrorProperties
	Records
}

func (o *RunReportOutput) decode(body io.ReadCloser) error { return unmarshalJSON(body, &o) }

// RunReport sends a request to POST /v1/reports/{reportId}/run?tableId={tableId}.
// See https://developer.quickbase.com/operation/runReport
func (c *Client) RunReport(input *RunReportInput) (output *RunReportOutput, err error) {
	input.c = c
	input.u = c.URL + "/reports/" + url.PathEscape(input.ReportID) + "/run?tableId=" + url.QueryEscape(input.TableID)
	output = &RunReportOutput{}
	err = c.Do(input, output)
	return
}
