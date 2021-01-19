package qbclient

import (
	"io"
	"net/http"
	"net/url"
)

// ListRelationshipsInput models the input sent to GET /v1/tables/{tableId}/relationships.
// See https://developer.quickbase.com/operation/getRelationships
type ListRelationshipsInput struct {
	c *Client
	u string

	ChildTableID string `json:"-" validate:"required" cliutil:"option=child-table-id"`
}

func (i *ListRelationshipsInput) url() string                  { return i.u }
func (i *ListRelationshipsInput) method() string               { return http.MethodGet }
func (i *ListRelationshipsInput) addHeaders(req *http.Request) { addHeadersJSON(req, i.c) }
func (i *ListRelationshipsInput) encode() ([]byte, error)      { return marshalJSON(i) }

// ListRelationshipsOutput models the output returned by GET /v1/tables/{tableId}/relationships.
// See https://developer.quickbase.com/operation/getRelationships
type ListRelationshipsOutput struct {
	ErrorProperties

	Metadata      *ListRelationshipsOutputMetadata `json:"metadata,omitempty"`
	Relationships []*Relationship                  `json:"relationships,omitempty"`
}

func (o *ListRelationshipsOutput) decode(body io.ReadCloser) error { return unmarshalJSON(body, &o) }

// ListRelationshipsOutputMetadata models the metadata property.
type ListRelationshipsOutputMetadata struct {
	NumberOfRelationships int `json:"numRelationships,omitempty"`
	Skip                  int `json:"skip,omitempty"`
	TotalRelationships    int `json:"totalRelationships,omitempty"`
}

// ListRelationships sends a request to GET /v1/tables/{tableId}/relationships.
// See https://developer.quickbase.com/operation/getRelationships
func (c *Client) ListRelationships(input *ListRelationshipsInput) (output *ListRelationshipsOutput, err error) {
	input.c = c
	input.u = c.URL + "/tables/" + url.PathEscape(input.ChildTableID) + "/relationships"
	output = &ListRelationshipsOutput{}
	err = c.Do(input, output)
	return
}

// ListRelationshipsByTableID sends a request to GET /v1/tables/{tableId}/relationships
// and gets a relationship by table ID.
// See https://developer.quickbase.com/operation/getTable
func (c *Client) ListRelationshipsByTableID(id string) (*ListRelationshipsOutput, error) {
	return c.ListRelationships(&ListRelationshipsInput{ChildTableID: id})
}

// CreateRelationshipInput models the input sent to POST /v1/tables/{tableId}/relationship.
// See https://developer.quickbase.com/operation/createRelationship
type CreateRelationshipInput struct {
	c *Client
	u string

	ChildTableID    string                                  `json:"-" validate:"required" cliutil:"option=child-table-id"`
	ParentTableID   string                                  `json:"parentTableId,omitempty" validate:"required" cliutil:"option=parent-table-id"`
	ForeignKeyField *CreateRelationshipInputForeignKeyField `json:"foreignKeyField,omitempty"`
	LookupFieldIDs  []int                                   `json:"lookupFieldIds,omitempty" cliutil:"option=lookup-field-ids"`
	SummaryFields   []*RelationshipSummaryField             `json:"summaryFields,omitempty"`
}

func (i *CreateRelationshipInput) url() string                  { return i.u }
func (i *CreateRelationshipInput) method() string               { return http.MethodPost }
func (i *CreateRelationshipInput) addHeaders(req *http.Request) { addHeadersJSON(req, i.c) }
func (i *CreateRelationshipInput) encode() ([]byte, error)      { return marshalJSON(i) }

// CreateRelationshipInputForeignKeyField models the summaryFields property.
type CreateRelationshipInputForeignKeyField struct {
	Label string `json:"label,omitempty" cliutil:"option=foreign-key-label"`
}

// CreateRelationshipOutput models the output returned by POST /v1/tables/{tableId}/relationship.
// See https://developer.quickbase.com/operation/createRelationship
type CreateRelationshipOutput struct {
	ErrorProperties
	Relationship
}

func (o *CreateRelationshipOutput) decode(body io.ReadCloser) error {
	return unmarshalJSON(body, &o)
}

// CreateRelationship sends a request to POST /v1/tables/{tableId}/relationship.
// See https://developer.quickbase.com/operation/createRelationship
func (c *Client) CreateRelationship(input *CreateRelationshipInput) (output *CreateRelationshipOutput, err error) {
	input.c = c
	input.u = c.URL + "/tables/" + url.PathEscape(input.ChildTableID) + "/relationship"
	output = &CreateRelationshipOutput{}
	err = c.Do(input, output)
	return
}

// UpdateRelationshipInput models the input sent to POST /v1/tables/{tableId}/relationship/{relationshipId}.
// See https://developer.quickbase.com/operation/updateRelationship
type UpdateRelationshipInput struct {
	c *Client
	u string

	ChildTableID   string                      `json:"-" validate:"required" cliutil:"option=child-table-id"`
	RelationshipID int                         `json:"-" validate:"required" cliutil:"option=relationship-id"`
	LookupFieldIDs []int                       `json:"lookupFieldIds,omitempty" cliutil:"option=lookup-field-ids"`
	SummaryFields  []*RelationshipSummaryField `json:"summaryFields,omitempty"`
}

func (i *UpdateRelationshipInput) url() string                  { return i.u }
func (i *UpdateRelationshipInput) method() string               { return http.MethodPost }
func (i *UpdateRelationshipInput) addHeaders(req *http.Request) { addHeadersJSON(req, i.c) }
func (i *UpdateRelationshipInput) encode() ([]byte, error)      { return marshalJSON(i) }

// UpdateRelationshipOutput models the output returned by POST /v1/tables/{tableId}/relationship/{relationshipId}.
// See https://developer.quickbase.com/operation/updateRelationship
type UpdateRelationshipOutput struct {
	ErrorProperties
	Relationship
}

func (o *UpdateRelationshipOutput) decode(body io.ReadCloser) error {
	return unmarshalJSON(body, &o)
}

// UpdateRelationship sends a request to POST /v1/tables/{tableId}/relationship/{relationshipId}.
// See https://developer.quickbase.com/operation/updateRelationship
func (c *Client) UpdateRelationship(input *UpdateRelationshipInput) (output *UpdateRelationshipOutput, err error) {
	input.c = c
	input.u = c.URL + relationshipPath(input.ChildTableID, input.RelationshipID)
	output = &UpdateRelationshipOutput{}
	err = c.Do(input, output)
	return
}

// DeleteRelationshipInput models the input sent to DELETE /v1/tables/{tableId}/relationship/{relationshipId}
// See https://developer.quickbase.com/operation/deleteRelationship
type DeleteRelationshipInput struct {
	c *Client
	u string

	ChildTableID   string `json:"-" validate:"required" cliutil:"option=child-table-id"`
	RelationshipID int    `json:"-" validate:"required" cliutil:"option=relationship-id"`
}

func (i *DeleteRelationshipInput) url() string                  { return i.u }
func (i *DeleteRelationshipInput) method() string               { return http.MethodDelete }
func (i *DeleteRelationshipInput) addHeaders(req *http.Request) { addHeadersJSON(req, i.c) }
func (i *DeleteRelationshipInput) encode() ([]byte, error)      { return marshalJSON(i) }

// DeleteRelationshipOutput models the output returned by DELETE /v1/tables/{tableId}/relationship/{relationshipId}
// See https://developer.quickbase.com/operation/deleteRelationship
type DeleteRelationshipOutput struct {
	ErrorProperties

	RelationshipID int `json:"relationshipId,omitempty"`
}

func (o *DeleteRelationshipOutput) decode(body io.ReadCloser) error {
	return unmarshalJSON(body, &o)
}

// DeleteRelationship sends a request to DELETE /v1/tables/{tableId}/relationship/{relationshipId}
// See https://developer.quickbase.com/operation/deleteRelationship
func (c *Client) DeleteRelationship(input *DeleteRelationshipInput) (output *DeleteRelationshipOutput, err error) {
	input.c = c
	input.u = c.URL + relationshipPath(input.ChildTableID, input.RelationshipID)
	output = &DeleteRelationshipOutput{}
	err = c.Do(input, output)
	return
}
