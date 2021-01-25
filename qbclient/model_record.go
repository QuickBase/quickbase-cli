package qbclient

import (
	"encoding/json"
	"fmt"
)

// Record models a record in Quick Base.
type Record struct {
	Fields map[int]*Value
}

// SetValue sets a value for a field.
func (r *Record) SetValue(fid int, val *Value) {
	if len(r.Fields) == 0 {
		r.Fields = make(map[int]*Value)
	}
	r.Fields[fid] = val
}

// UnmarshalJSON implements json.UnmarshalJSON by using the field type to
// decode the "value" parameter into the appropriate data type.
//
// TODO DRY up the code, see below.
func (output *QueryRecordsOutput) UnmarshalJSON(b []byte) (err error) {

	// Unmarshal the json into our parsing struct.
	var v parseQueryRecordsOutput
	if err = json.Unmarshal(b, &v); err != nil {
		return
	}

	// Set the parsed value for everything but the "data" property.
	output.Message = v.Message
	output.Description = v.Description
	output.Fields = v.Fields
	output.Metadata = v.Metadata

	// Build a mapping of field IDs to Quick Base field type.
	tmap := make(map[int]string, len(v.Fields))
	for _, fd := range v.Fields {
		tmap[fd.FieldID] = fd.Type
	}

	// Parse the field values now that we have the Quick Base field types.
	output.Data = make([]map[int]*RecordsData, len(v.Data))
	for i, record := range v.Data {

		data := make(map[int]*RecordsData, len(record))
		for fid, field := range record {

			// Get the Quick Base field type from the fid.
			ftype, ok := tmap[fid]
			if !ok {
				err = fmt.Errorf("fid %v, %w", fid, fmt.Errorf("%s: %w", ftype, ErrInvalidType))
				return
			}

			// Unmarshal the field based on its Quick Base type.
			var val *Value
			if val, err = unmarshalField(fid, ftype, field.Value); err != nil {
				return
			}

			data[fid] = &RecordsData{Value: val}
		}

		output.Data[i] = data
	}

	return
}

// UnmarshalJSON implements json.UnmarshalJSON by using the field type to
// decode the "value" parameter into the appropriate data type.
//
// TODO DRY up the code, see above.
func (output *RunReportOutput) UnmarshalJSON(b []byte) (err error) {

	// Unmarshal the json into our parsing struct.
	var v parseQueryRecordsOutput
	if err = json.Unmarshal(b, &v); err != nil {
		return
	}

	// Set the parsed value for everything but the "data" property.
	output.Message = v.Message
	output.Description = v.Description
	output.Fields = v.Fields
	output.Metadata = v.Metadata

	// Build a mapping of field IDs to Quick Base field type.
	tmap := make(map[int]string, len(v.Fields))
	for _, fd := range v.Fields {
		tmap[fd.FieldID] = fd.Type
	}

	// Parse the field values now that we have the Quick Base field types.
	output.Data = make([]map[int]*RecordsData, len(v.Data))
	for i, record := range v.Data {

		data := make(map[int]*RecordsData, len(record))
		for fid, field := range record {

			// Get the Quick Base field type from the fid.
			ftype, ok := tmap[fid]
			if !ok {
				err = fmt.Errorf("fid %v, %w", fid, fmt.Errorf("%s: %w", ftype, ErrInvalidType))
				return
			}

			// Unmarshal the field based on its Quick Base type.
			var val *Value
			if val, err = unmarshalField(fid, ftype, field.Value); err != nil {
				return
			}

			data[fid] = &RecordsData{Value: val}
		}

		output.Data[i] = data
	}

	return
}

type parseQueryRecordsOutput struct {
	ErrorProperties

	Fields   []*RecordsField  `json:"fields,omitempty"`
	Metadata *RecordsMetadata `json:"metadata,omitempty"`
	Data     []map[int]struct {
		Value *json.RawMessage `json:"value"`
	} `json:"data,omitempty"`
}

// SetRecords sets the records to insert.
//
// This function converts the a Records slice and sets it as the
// InsertRecordsInput.Data property.
func (i *InsertRecordsInput) SetRecords(records []*Record) {
	i.Data = make([]map[int]*InsertRecordsInputData, len(records))
	for n, r := range records {
		data := make(map[int]*InsertRecordsInputData)
		for fid, val := range r.Fields {
			data[fid] = &InsertRecordsInputData{Value: val}
		}
		i.Data[n] = data
	}
}

// Records models output that returns records.
type Records struct {
	Data     []map[int]*RecordsData `json:"data,omitempty"`
	Fields   []*RecordsField        `json:"fields,omitempty"`
	Metadata *RecordsMetadata       `json:"metadata,omitempty"`
}

// RecordsData models objects in the data array.
type RecordsData struct {
	Value *Value `json:"value"`
}

// RecordsField models objects in the fields array.
type RecordsField struct {
	FieldID int    `json:"id"`
	Label   string `json:"label"`
	Type    string `json:"type"`
}

// RecordsMetadata models the metadata object.
type RecordsMetadata struct {
	TotalRecords int `json:"totalRecords"`
	NumRecords   int `json:"numRecords"`
	NumFields    int `json:"numFields"`
	Skip         int `json:"skip"`
	Top          int `json:"top,omitempty"`
}
