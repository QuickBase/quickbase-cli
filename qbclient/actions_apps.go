package qbclient

import (
	"io"
	"net/http"
)

// ListAppsInput models the XML API request sent to API_GrantedDBs.
// See https://help.quickbase.com/api-guide/granteddbs.html
type ListAppsInput struct {
	XMLRequestParameters
	XMLCredentialParameters

	c *Client
	u string

	AdminOnly          bool `xml:"adminOnly,omitempty" cliutil:"option=admin-only"`
	ExcludeParents     bool `xml:"excludeparents,int,omitempty" cliutil:"option=exclude-parents"`
	IncludeAncestors   bool `xml:"includeancestors,int,omitempty" cliutil:"option=include-ancestors"`
	RealmAppsOnly      bool `xml:"realmAppsOnly,omitempty" cliutil:"option=realm-apps-only"`
	WithEmbeddedTables bool `xml:"withembeddedtables,int" cliutil:"option=with-embedded-tables"`
}

func (i *ListAppsInput) method() string               { return http.MethodPost }
func (i *ListAppsInput) url() string                  { return i.u }
func (i *ListAppsInput) addHeaders(req *http.Request) { addHeadersXML(req, i.c, "API_GrantedDBs") }
func (i *ListAppsInput) encode() ([]byte, error)      { return marshalXML(i, i.c) }

// ListAppsOutput models the XML API response returned by API_GrantedDBs.
// See https://help.quickbase.com/api-guide/granteddbs.html
type ListAppsOutput struct {
	XMLResponseParameters

	Databases []*ListAppsOutputDatabases `xml:"databases>dbinfo" json:"apps,omitempty"`
}

// ListAppsOutputDatabases modesl the databases propertie.
type ListAppsOutputDatabases struct {
	AncestorAppID       string `xml:"ancestorappid,omitempty" json:"ancestorAppId,omitempty"`
	ID                  string `xml:"dbid" json:"appId"`
	Name                string `xml:"dbname" json:"name"`
	OldestAncestorAppID string `xml:"oldestancestorappid,omitempty"  json:"oldAncestorAppId,omitempty"`
}

func (o *ListAppsOutput) decode(body io.ReadCloser) error { return unmarshalXML(body, o) }

// ListApps sends an XML API request to API_GrantedDBs.
// See https://help.quickbase.com/api-guide/granteddbs.html
func (c *Client) ListApps(input *ListAppsInput) (output *ListAppsOutput, err error) {
	input.c = c
	input.u = "https://" + c.ReamlHostname + "/db/main"
	output = &ListAppsOutput{}
	err = c.Do(input, output)
	return
}
