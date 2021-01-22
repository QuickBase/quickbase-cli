package qbclient

import (
	"io"
	"net/http"
	"net/url"

	"github.com/QuickBase/quickbase-cli/qberrors"
)

// CreatePageInput models the XML API request sent to API_AddReplaceDBPage
// See https://help.quickbase.com/api-guide/index.html#add_replace_dbpage.html
type CreatePageInput struct {
	XMLRequestParameters
	XMLCredentialParameters

	c *Client
	u string

	AppID string               `xml:"-" validate:"required" cliutil:"option=app-id"`
	Body  *CreatePageInputBody `xml:"pagebody"`
	Name  string               `xml:"pagename" validate:"required" cliutil:"option=page-name"`
	Type  int                  `xml:"pagetype" validate:"required" cliutil:"option=page-type default=1"`
}

func (i *CreatePageInput) method() string { return http.MethodPost }
func (i *CreatePageInput) url() string    { return i.u }
func (i *CreatePageInput) addHeaders(req *http.Request) {
	addHeadersXML(req, i.c, "API_AddReplaceDBPage")
}
func (i *CreatePageInput) encode() ([]byte, error) { return marshalXML(i, i.c) }

// CreatePageInputBody models the pagebody element.
type CreatePageInputBody struct {
	Data string `xml:",cdata" validate:"required" cliutil:"option=page-body func=ioreader"`
}

// CreatePageOutput models the XML API response returned by API_AddReplaceDBPage
// See https://help.quickbase.com/api-guide/index.html#add_replace_dbpage.html
type CreatePageOutput struct {
	XMLResponseParameters

	PageID int `xml:"pageID"`
}

func (o *CreatePageOutput) decode(body io.ReadCloser) error { return unmarshalXML(body, o) }

// CreatePage sends an XML API request to API_AddReplaceDBPage.
// See https://help.quickbase.com/api-guide/index.html#add_replace_dbpage.html
func (c *Client) CreatePage(input *CreatePageInput) (output *CreatePageOutput, err error) {
	input.c = c
	input.u = "https://" + url.PathEscape(c.ReamlHostname) + "/db/" + url.PathEscape(input.AppID)
	output = &CreatePageOutput{}
	err = c.Do(input, output)
	return
}

// GetPageInput models the XML API request sent to API_GetDBPage
// See https://help.quickbase.com/api-guide/index.html#get_db_page.html
type GetPageInput struct {
	XMLRequestParameters
	XMLCredentialParameters

	c *Client
	u string

	AppID  string `xml:"-" validate:"required" cliutil:"option=app-id"`
	PageID string `xml:"pageID" validate:"required" cliutil:"option=page-id"`
}

func (i *GetPageInput) method() string               { return http.MethodPost }
func (i *GetPageInput) url() string                  { return i.u }
func (i *GetPageInput) addHeaders(req *http.Request) { addHeadersXML(req, i.c, "API_GetDBPage") }
func (i *GetPageInput) encode() ([]byte, error)      { return marshalXML(i, i.c) }

// GetPageOutput models the XML API response returned by API_GetDBPage
// See https://help.quickbase.com/api-guide/index.html#get_db_page.html
type GetPageOutput struct {
	XMLResponseParameters

	Body string `xml:"pagebody"`
}

func (o *GetPageOutput) decode(body io.ReadCloser) error { return unmarshalXML(body, o) }

// GetPage makes an API_GetDBPage call.
// See https://help.quickbase.com/api-guide/index.html#get_db_page.html
func (c *Client) GetPage(input *GetPageInput) (output *GetPageOutput, err error) {
	input.c = c
	input.u = "https://" + url.PathEscape(c.ReamlHostname) + "/db/" + url.PathEscape(input.AppID)
	output = &GetPageOutput{}
	err = c.Do(input, output)
	return
}

// UpdatePageInput models the XML API request sent to API_AddReplaceDBPage
// See https://help.quickbase.com/api-guide/index.html#add_replace_dbpage.html
type UpdatePageInput struct {
	XMLRequestParameters
	XMLCredentialParameters

	c *Client
	u string

	AppID  string               `xml:"-" validate:"required" cliutil:"option=app-id"`
	PageID int                  `xml:"pageid" cliutil:"option=page-id"`
	Name   string               `xml:"pagename" cliutil:"option=page-name"`
	Body   *UpdatePageInputBody `xml:"pagebody"`
}

func (i *UpdatePageInput) method() string { return http.MethodPost }
func (i *UpdatePageInput) url() string    { return i.u }
func (i *UpdatePageInput) addHeaders(req *http.Request) {
	addHeadersXML(req, i.c, "API_AddReplaceDBPage")
}
func (i *UpdatePageInput) encode() ([]byte, error) { return marshalXML(i, i.c) }

// UpdatePageInputBody models the pagebody element.
type UpdatePageInputBody struct {
	Data string `xml:",cdata" validate:"required" cliutil:"option=page-body func=ioreader"`
}

// UpdatePageOutput models the XML API response returned by API_AddReplaceDBPage
// See https://help.quickbase.com/api-guide/index.html#add_replace_dbpage.html
type UpdatePageOutput struct {
	XMLResponseParameters

	PageID int `xml:"pageID"`
}

func (o *UpdatePageOutput) decode(body io.ReadCloser) error { return unmarshalXML(body, o) }

// UpdatePage sends an XML API request to API_AddReplaceDBPage.
// See https://help.quickbase.com/api-guide/index.html#add_replace_dbpage.html
func (c *Client) UpdatePage(input *UpdatePageInput) (output *UpdatePageOutput, err error) {
	input.c = c
	input.u = "https://" + url.PathEscape(c.ReamlHostname) + "/db/" + url.PathEscape(input.AppID)
	output = &UpdatePageOutput{}
	if input.PageID != 0 || input.Name != "" {
		err = c.Do(input, output)
	} else {
		err = qberrors.Client(nil).Safef(qberrors.BadRequest, "ID or name required")
	}
	return
}
