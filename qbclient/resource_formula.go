package qbclient

import (
	"io"
	"net/http"
)

// RunFormulaInput models the input sent to POST /v1/formula/run.
// See https://developer.quickbase.com/operation/runFormula
type RunFormulaInput struct {
	c *Client
	u string

	From     string `json:"from" validate:"required" cliutil:"option=from"`
	RecordID int    `json:"rid" validate:"required" cliutil:"option=record-id"`
	Formula  string `json:"formula" validate:"required" cliutil:"option=formula func=stdin"`
}

func (i *RunFormulaInput) url() string                  { return i.u }
func (i *RunFormulaInput) method() string               { return http.MethodPost }
func (i *RunFormulaInput) addHeaders(req *http.Request) { addHeadersJSON(req, i.c) }
func (i *RunFormulaInput) encode() ([]byte, error)      { return marshalJSON(i) }

// RunFormulaOutput models the output returned by POST /v1/formula/run.
// See https://developer.quickbase.com/operation/runFormula
type RunFormulaOutput struct {
	ErrorProperties

	Result string `json:"result,omitempty"`
}

func (o *RunFormulaOutput) decode(body io.ReadCloser) error { return unmarshalJSON(body, &o) }

// RunFormula sends a request to POST /v1/formula/run.
// See https://developer.quickbase.com/operation/runFormula
func (c *Client) RunFormula(input *RunFormulaInput) (output *RunFormulaOutput, err error) {
	input.c = c
	input.u = c.URL + "/formula/run"
	output = &RunFormulaOutput{}
	err = c.Do(input, output)
	return
}
