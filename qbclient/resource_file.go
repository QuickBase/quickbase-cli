package qbclient

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
)

// DeleteFileInput models the input sent to DELETE /v1/files/{tableId}/{recordId}/{fieldId}/{versionNumber}.
// See https://developer.quickbase.com/operation/deleteFile
type DeleteFileInput struct {
	c *Client
	u string

	TableID  string `json:"-" validate:"required" cliutil:"option=table-id"`
	RecordID int    `json:"-" validate:"required" cliutil:"option=record-id"`
	FieldID  int    `json:"-" validate:"required" cliutil:"option=field-id"`
	Version  int    `json:"-" validate:"required" cliutil:"option=version"`
}

func (i *DeleteFileInput) url() string                  { return i.u }
func (i *DeleteFileInput) method() string               { return http.MethodDelete }
func (i *DeleteFileInput) addHeaders(req *http.Request) { addHeadersJSON(req, i.c) }
func (i *DeleteFileInput) encode() ([]byte, error)      { return marshalJSON(i) }

// DeleteFileOutput models the output returned by DELETE /v1/files/{tableId}/{recordId}/{fieldId}/{versionNumber}.
// See https://developer.quickbase.com/operation/deleteFile
type DeleteFileOutput struct {
	ErrorProperties

	Version  int                      `json:"versionNumber,omitempty"`
	FileName string                   `json:"fileName,omitempty"`
	Uploaded *Timestamp               `json:"uploaded,omitempty"`
	Creator  *DeleteFileOutputCreator `json:"creator,omitempty"`
}

func (o *DeleteFileOutput) decode(body io.ReadCloser) error { return unmarshalJSON(body, &o) }

// DeleteFileOutputCreator models the creator object.
type DeleteFileOutputCreator struct {
	Email  string `json:"email,omitempty"`
	UserID string `json:"id,omitempty"`
	Name   string `json:"name,omitempty"`
}

// DeleteFile sends a request to DELETE /v1/files/{tableId}/{recordId}/{fieldId}/{versionNumber}.
// See https://developer.quickbase.com/operation/deleteFile
func (c *Client) DeleteFile(input *DeleteFileInput) (output *DeleteFileOutput, err error) {
	input.c = c
	input.u = fmt.Sprintf("%s/files/%s/%v/%v/%v", c.URL, url.PathEscape(input.TableID), input.RecordID, input.FieldID, input.Version)
	output = &DeleteFileOutput{}
	err = c.Do(input, output)
	return
}
