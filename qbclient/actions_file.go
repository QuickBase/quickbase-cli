package qbclient

import (
	"encoding/base64"
	"io"
	"net/http"
	"net/url"
)

// CreateFileInput models the XML API request sent to API_UploadFile
// See https://help.quickbase.com/api-guide/index.html#uploadfile.html
type CreateFileInput struct {
	XMLRequestParameters
	XMLCredentialParameters

	c *Client
	u string

	TableID  string                  `xml:"-" validate:"required" cliutil:"option=table-id"`
	Fields   []*CreateFileInputField `xml:"field"`
	RecordID int                     `xml:"rid" validate:"required" cliutil:"option=record-id"`
}

func (i *CreateFileInput) method() string               { return http.MethodPost }
func (i *CreateFileInput) url() string                  { return i.u }
func (i *CreateFileInput) addHeaders(req *http.Request) { addHeadersXML(req, i.c, "API_UploadFile") }
func (i *CreateFileInput) encode() ([]byte, error)      { return marshalXML(i, i.c) }

// CreateFileInputField models the field element.
type CreateFileInputField struct {
	FieldID  int    `xml:"fid,attr" validate:"required" cliutil:"option=field-id"`
	FileData string `xml:",chardata" validate:"required" cliutil:"option=file-data func=ioreader"`
	Name     string `xml:"filename,attr" validate:"required" cliutil:"option=file-name"`
}

// CreateFileOutput models the XML API response returned by API_UploadFile
// See https://help.quickbase.com/api-guide/index.html#uploadfile.html
type CreateFileOutput struct {
	XMLResponseParameters

	Fields []*CreateFileOutputField `xml:"file_fields>field"`
}

func (o *CreateFileOutput) decode(body io.ReadCloser) error { return unmarshalXML(body, o) }

// CreateFileOutputField models the file_fields element.
type CreateFileOutputField struct {
	FileID int    `xml:"id,attr"`
	URL    string `xml:"url"`
}

// CreateFile makes an API_UploadFile call.
// See https://help.quickbase.com/api-guide/index.html#uploadfile.html
func (c *Client) CreateFile(input *CreateFileInput) (output *CreateFileOutput, err error) {
	input.c = c
	input.u = "https://" + c.ReamlHostname + "/db/" + url.PathEscape(input.TableID)

	for _, f := range input.Fields {
		f.FileData = base64.StdEncoding.EncodeToString([]byte(f.FileData))
	}

	output = &CreateFileOutput{}
	err = c.Do(input, output)
	return
}
