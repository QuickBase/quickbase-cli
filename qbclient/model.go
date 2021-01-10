package qbclient

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/araddon/dateparse"
)

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

//
// Developer-friendly models for records and fields. The custom marshaler and
// unmarshaler of the Value is where serialization and deserialization
// happen.
//

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

// Value models the value of fields in Quick Base. This struct effectively
// handles the Quick base field type / Golang type transformations.
type Value struct {
	Bool        bool
	Duration    time.Duration
	Float64     float64
	String      string
	StringSlice []string
	Time        time.Time
	URL         *url.URL
	User        *User
	UserSlice   []*User

	QuickBaseType string
}

// NewRecordIDValue returns a new Value of the FieldRecordID type.
func NewRecordIDValue(val float64) *Value {
	return &Value{Float64: val, QuickBaseType: FieldRecordID}
}

// NewTextValue returns a new Value of the FieldText type.
func NewTextValue(val string) *Value {
	return &Value{String: val, QuickBaseType: FieldText}
}

// NewTextMultiLineValue returns a new Value of the FieldTextMultiLine type.
func NewTextMultiLineValue(val string) *Value {
	return &Value{String: val, QuickBaseType: FieldTextMultiLine}
}

// NewTextMultipleChoiceValue returns a new Value of the FieldTextMultipleChoice type.
func NewTextMultipleChoiceValue(val string) *Value {
	return &Value{String: val, QuickBaseType: FieldTextMultipleChoice}
}

// NewRichTextValue returns a new Value of the FieldRichText type.
func NewRichTextValue(val string) *Value {
	return &Value{String: val, QuickBaseType: FieldRichText}
}

// NewMultiSelectTextValue returns a new Value of the FieldMultiSelectText type.
func NewMultiSelectTextValue(val []string) *Value {
	return &Value{StringSlice: val, QuickBaseType: FieldMultiSelectText}
}

// NewMultiSelectTextValueFromString returns a new Value of the FieldMultiSelectText
// type given a string with a comma-separated list of values.
func NewMultiSelectTextValueFromString(val string) (v *Value, err error) {
	var ss []string
	if ss, err = ParseList(val); err == nil {
		v = NewMultiSelectTextValue(ss)
	}
	return
}

// NewNumericValue returns a new Value of the FieldNumeric type.
func NewNumericValue(val float64) *Value {
	return &Value{Float64: val, QuickBaseType: FieldNumeric}
}

// NewNumericValueFromString returns a new Value of the FieldNumeric type
// given a string.
func NewNumericValueFromString(val string) (*Value, error) {
	return parseStringToNumericValue(val, FieldNumeric)
}

// NewNumericCurrencyValue returns a new Value of the FieldNumericCurrency type.
func NewNumericCurrencyValue(val float64) *Value {
	return &Value{Float64: val, QuickBaseType: FieldNumericCurrency}
}

// NewNumericCurrencyValueFromString returns a new Value of the
// FieldNumericCurrency type given a string
func NewNumericCurrencyValueFromString(val string) (*Value, error) {
	return parseStringToNumericValue(val, FieldNumericCurrency)
}

// NewNumericPercentValue returns a new Value of the FieldNumericPercent type.
func NewNumericPercentValue(val float64) *Value {
	return &Value{Float64: val, QuickBaseType: FieldNumericPercent}
}

// NewNumericPercentValueFromString returns a new Value of the FieldNumericPercent
// type given a string.
func NewNumericPercentValueFromString(val string) (*Value, error) {
	return parseStringToNumericValue(val, FieldNumericPercent)
}

// NewNumericRatingValue returns a new Value of the FieldNumericRating type.
func NewNumericRatingValue(val float64) *Value {
	return &Value{Float64: val, QuickBaseType: FieldNumericRating}
}

// NewNumericRatingValueFromString returns a new Value of the FieldNumericRating
// type given a string.
func NewNumericRatingValueFromString(val string) (*Value, error) {
	return parseStringToNumericValue(val, FieldNumericRating)
}

// NewDateValue returns a new Value of the FieldDate type.
func NewDateValue(val time.Time) *Value {
	return &Value{Time: val, QuickBaseType: FieldDate}
}

// NewDateValueFromString returns a new Value of the FieldDate type,
// parsing the passed string into a time.Time.
func NewDateValueFromString(val string) (*Value, error) {
	return parseTimeToValue(val, FieldDate, dateparse.ParseAny)
}

// NewDateTimeValue returns a new Value of the FieldDateTime type.
func NewDateTimeValue(val time.Time) *Value {
	return &Value{Time: val, QuickBaseType: FieldDateTime}
}

// NewDateTimeValueFromString returns a new Value of the FieldDate type,
// parsing the passed string into a time.Time.
func NewDateTimeValueFromString(val string) (*Value, error) {
	return parseTimeToValue(val, FieldDateTime, dateparse.ParseLocal)
}

// NewTimeOfDayValue returns a new Value of the FieldTimeOfDay type.
func NewTimeOfDayValue(val time.Time) *Value {
	return &Value{Time: val, QuickBaseType: FieldTimeOfDay}
}

// NewTimeOfDayValueFromString returns a new Value of the FieldDate type,
// parsing the passed string into a time.Time.
func NewTimeOfDayValueFromString(val string) (*Value, error) {
	return parseTimeToValue("3/19/1982 "+val, FieldTimeOfDay, dateparse.ParseAny)
}

// NewDurationValue returns a new Value of the FieldDuration type.
func NewDurationValue(val time.Duration) *Value {
	return &Value{Duration: val, QuickBaseType: FieldDuration}
}

// NewDurationValueFromFloat64 returns a new Value of the FieldDuration type,
// converting the passed float64 into a duration. We assume that the float64 is
// the duration in milliseconds.
func NewDurationValueFromFloat64(val float64) *Value {
	return NewDurationValue(time.Millisecond * time.Duration(val))
}

// NewDurationValueFromString returns a new Value of the FieldDuration type
// given a passed string.
func NewDurationValueFromString(val string) (v *Value, err error) {
	var d time.Duration
	if d, err = time.ParseDuration(val); err == nil {
		v = NewDurationValue(d)
	}
	return
}

// NewCheckboxValue returns a new Value of the FieldCheckbox type.
func NewCheckboxValue(val bool) *Value {
	return &Value{Bool: val, QuickBaseType: FieldCheckbox}
}

// NewCheckboxValueFromString returns a new Value of the FieldCheckbox type
// given a passed string.
func NewCheckboxValueFromString(val string) (v *Value, err error) {
	var b bool
	if b, err = strconv.ParseBool(val); err != nil {
		v = NewCheckboxValue(b)
	}
	return
}

// NewAddressValue returns a new Value of the FieldAddress type.
func NewAddressValue(val string) *Value {
	return &Value{String: val, QuickBaseType: FieldAddress}
}

// NewAddressStreet1Value returns a new Value of the FieldAddressStreet1 type.
func NewAddressStreet1Value(val string) *Value {
	return &Value{String: val, QuickBaseType: FieldAddressStreet1}
}

// NewAddressStreet2Value returns a new Value of the FieldAddressStreet2 type.
func NewAddressStreet2Value(val string) *Value {
	return &Value{String: val, QuickBaseType: FieldAddressStreet2}
}

// NewAddressCityValue returns a new Value of the FieldAddressCity type.
func NewAddressCityValue(val string) *Value {
	return &Value{String: val, QuickBaseType: FieldAddressCity}
}

// NewAddressStateRegionValue returns a new Value of the FieldAddressStateRegion type.
func NewAddressStateRegionValue(val string) *Value {
	return &Value{String: val, QuickBaseType: FieldAddressStateRegion}
}

// NewAddressPostalCodeValue returns a new Value of the FieldAddressPostalCode type.
func NewAddressPostalCodeValue(val string) *Value {
	return &Value{String: val, QuickBaseType: FieldAddressPostalCode}
}

// NewAddressCountryValue returns a new Value of the FieldAddressCountry type.
func NewAddressCountryValue(val string) *Value {
	return &Value{String: val, QuickBaseType: FieldAddressCountry}
}

// NewPhoneNumberValue returns a new Value of the FieldPhoneNumber type.
func NewPhoneNumberValue(val string) *Value {
	return &Value{String: val, QuickBaseType: FieldPhoneNumber}
}

// NewEmailAddressValue returns a new Value of the FieldEmailAddress type.
func NewEmailAddressValue(val string) *Value {
	return &Value{String: val, QuickBaseType: FieldEmailAddress}
}

// NewUserValue returns a new Value of the FieldUser type.
func NewUserValue(val *User) *Value {
	return &Value{User: val, QuickBaseType: FieldUser}
}

// NewUserValueFromString returns a new Value of the FieldUser type given a
// passed string.
func NewUserValueFromString(val string) *Value {
	return NewUserValue(&User{ID: val})
}

// NewListUserValue returns a new Value of the FieldUserList type.
func NewListUserValue(val []*User) *Value {
	return &Value{UserSlice: val, QuickBaseType: FieldUserList}
}

// NewListUserValueFromString returns a new Value of the FieldUserList type
// given a passed string.
func NewListUserValueFromString(val string) *Value {
	return &Value{UserSlice: []*User{}, QuickBaseType: FieldUserList}
}

// NewFileAttachmentValue returns a new Value of the FieldFileAttachment type.
func NewFileAttachmentValue(val string) *Value {
	return &Value{String: val, QuickBaseType: FieldFileAttachment}
}

// NewReportLinkValue returns a new Value of the FieldReportLink type.
func NewReportLinkValue(val string) *Value {
	return &Value{String: val, QuickBaseType: FieldReportLink}
}

// NewURLValue returns a new Value of the FieldURL type.
func NewURLValue(val *url.URL) *Value {
	return &Value{URL: val, QuickBaseType: FieldURL}
}

// NewURLValueFromString returns a new Value of the FieldURL type.
func NewURLValueFromString(val string) (v *Value, err error) {
	var u *url.URL
	if u, err = url.Parse(val); err != nil {
		v = &Value{URL: u, QuickBaseType: FieldURL}
	}
	return
}

// NewiCalendarValue returns a new Value of the FieldiCalendar type.
// Make this a URL?
func NewiCalendarValue(val string) *Value {
	return &Value{String: val, QuickBaseType: FieldiCalendar}
}

// NewvCardValue returns a new Value of the FieldvCard type.
// Make this a URL?
func NewvCardValue(val string) *Value {
	return &Value{String: val, QuickBaseType: FieldvCard}
}

// NewPredecessorValue returns a new Value of the FieldPredecessor type.
func NewPredecessorValue(val string) *Value {
	return &Value{String: val, QuickBaseType: FieldPredecessor}
}

type dateparseFn func(string, ...dateparse.ParserOption) (time.Time, error)

func parseTimeToValue(val, ftype string, fn dateparseFn) (v *Value, err error) {
	var t time.Time
	if val == "" {
		v = &Value{Time: t, QuickBaseType: ftype}
	} else if t, err = fn(val); err == nil {
		v = &Value{Time: t, QuickBaseType: ftype}
	}
	return
}

func parseStringToNumericValue(val, ftype string) (v *Value, err error) {
	var f float64
	if f, err = strconv.ParseFloat(val, 64); err == nil {
		v = &Value{Float64: f, QuickBaseType: ftype}
	}
	return
}

// NewValueFromString returns a new *Value from a string given the Quick Base
// field type.
func NewValueFromString(val, ftype string) (v *Value, err error) {
	switch ftype {

	case FieldText:
		// Also picks up:
		// FieldAddressStreet1, FieldAddressStreet2, FieldAddressCity,
		// FieldAddressStateRegion, FieldAddressPostalCode, and
		// FieldAddressCountry
		v = NewTextValue(val)

	case FieldTextMultiLine:
		v = NewTextMultiLineValue(val)

	case FieldTextMultipleChoice:
		v = NewTextMultipleChoiceValue(val)

	case FieldRichText:
		v = NewRichTextValue(val)

	case FieldMultiSelectText:
		v, err = NewMultiSelectTextValueFromString(val)

	case FieldNumeric:
		v, err = NewNumericValueFromString(val)

	case FieldNumericCurrency:
		v, err = NewNumericCurrencyValueFromString(val)

	case FieldNumericPercent:
		v, err = NewNumericCurrencyValueFromString(val)

	case FieldNumericRating:
		v, err = NewNumericPercentValueFromString(val)

	case FieldDate:
		v, err = NewDateValueFromString(val)

	case FieldDateTime:
		v, err = NewDateTimeValueFromString(val)

	case FieldTimeOfDay:
		v, err = NewTimeOfDayValueFromString(val)

	case FieldDuration:
		v, err = NewDurationValueFromString(val)

	case FieldCheckbox:
		v, err = NewCheckboxValueFromString(val)

	case FieldAddress:
		v = NewAddressValue(val)

	case FieldPhoneNumber:
		v = NewPhoneNumberValue(val)

	case FieldEmailAddress:
		v = NewEmailAddressValue(val)

	case FieldUser:
		v = NewUserValueFromString(val)

	case FieldUserList:
		v = NewListUserValueFromString(val)

	case FieldFileAttachment:
		v = NewFileAttachmentValue(val)

	case FieldReportLink:
		v = NewReportLinkValue(val)

	case FieldURL:
		v, err = NewURLValueFromString(val)

	case FieldiCalendar:
		v = NewiCalendarValue(val)

	case FieldvCard:
		v = NewvCardValue(val)

	case FieldPredecessor:
		v = NewPredecessorValue(val)

	default:
		err = fmt.Errorf("unsupported field type (%s)", ftype)
	}

	return
}

// MarshalJSON implements json.MarshalJSON and JSON encodes the value.
// TODO Marshal by Quick Base type instead, because we have to format dates differently.
func (v *Value) MarshalJSON() ([]byte, error) {
	switch v.QuickBaseType {

	case FieldRecordID:
		return json.Marshal(v.Float64)

	case FieldText, FieldTextMultiLine, FieldTextMultipleChoice, FieldRichText:
		// Also picks up:
		// FieldAddressStreet1, FieldAddressStreet2, FieldAddressCity,
		// FieldAddressStateRegion, FieldAddressPostalCode, and
		// FieldAddressCountry
		return json.Marshal(v.String)

	case FieldMultiSelectText:
		return json.Marshal(v.StringSlice)

	case FieldNumeric, FieldNumericCurrency, FieldNumericPercent, FieldNumericRating:
		return json.Marshal(v.Float64)

	case FieldDate:
		s := v.Time.UTC().Format(FormatDate)
		return json.Marshal(s)

	case FieldDateTime:
		s := v.Time.UTC().Format(FormatDateTime)
		return json.Marshal(s)

	case FieldTimeOfDay:
		s := v.Time.UTC().Format(FormatTimeOfDay)
		return json.Marshal(s)

	case FieldDuration:
		return json.Marshal(v.Duration.Milliseconds())

	case FieldCheckbox:
		return json.Marshal(v.Bool)

	case FieldAddress:
		return json.Marshal(v.String)

	case FieldPhoneNumber:
		return json.Marshal(v.String)

	case FieldEmailAddress:
		return json.Marshal(v.String)

	case FieldUser:
		return json.Marshal(v.User)

	case FieldUserList:
		return json.Marshal(v.UserSlice)

	// TODO revisit when re-added.
	case FieldFileAttachment:
		return json.Marshal(v.String)

	case FieldReportLink:
		return json.Marshal(v.String)

	case FieldURL:
		return json.Marshal(v.URL.String())

	// Deprecated
	case FieldiCalendar:
		return json.Marshal(v.String)

	// Deprecated
	case FieldvCard:
		return json.Marshal(v.String)

	case FieldPredecessor:
		return json.Marshal(v.String)

	default:
		return []byte(``), fmt.Errorf("unsupported field type (%s)", v.QuickBaseType)
	}
}

func unmarshalField(fid int, ftype string, data *json.RawMessage) (val *Value, err error) {
	switch ftype {

	case FieldRecordID:
		var v float64
		if data == nil {
			val = NewNumericValue(v)
		} else if err = json.Unmarshal(*data, &v); err == nil {
			val = NewNumericValue(v)
		}

	case FieldText:
		// Also picks up:
		// FieldAddressStreet1, FieldAddressStreet2, FieldAddressCity,
		// FieldAddressStateRegion, FieldAddressPostalCode, and
		// FieldAddressCountry
		var v string
		if data == nil {
			val = NewTextValue(v)
		} else if err = json.Unmarshal(*data, &v); err == nil {
			val = NewTextValue(v)
		}

	case FieldTextMultiLine:
		var v string
		if data == nil {
			val = NewTextMultiLineValue(v)
		} else if err = json.Unmarshal(*data, &v); err == nil {
			val = NewTextMultiLineValue(v)
		}

	case FieldTextMultipleChoice:
		var v string
		if data == nil {
			val = NewTextMultipleChoiceValue(v)
		} else if err = json.Unmarshal(*data, &v); err == nil {
			val = NewTextMultipleChoiceValue(v)
		}

	case FieldRichText:
		var v string
		if data == nil {
			val = NewRichTextValue(v)
		} else if err = json.Unmarshal(*data, &v); err == nil {
			val = NewRichTextValue(v)
		}

	case FieldMultiSelectText:
		var v []string
		if data == nil {
			val = NewMultiSelectTextValue(v)
		} else if err = json.Unmarshal(*data, &v); err == nil {
			val = NewMultiSelectTextValue(v)
		}

	case FieldNumeric, FieldNumericCurrency, FieldNumericPercent, FieldNumericRating:
		var v float64
		if data == nil {
			val = NewNumericValue(v)
		} else if err = json.Unmarshal(*data, &v); err == nil {
			val = NewNumericValue(v)
		}

	case FieldDate:
		var v string
		if data == nil {
			val, err = NewDateValueFromString(v)
		} else if err = json.Unmarshal(*data, &v); err == nil {
			val, err = NewDateValueFromString(v)
		}

	case FieldDateTime:
		var v string
		if data == nil {
			val, err = NewDateTimeValueFromString(v)
		} else if err = json.Unmarshal(*data, &v); err == nil {
			val, err = NewDateTimeValueFromString(v)
		}

	case FieldTimeOfDay:
		var v string
		if data == nil {
			val, err = NewTimeOfDayValueFromString(v)
		} else if err = json.Unmarshal(*data, &v); err == nil {
			val, err = NewTimeOfDayValueFromString(v)
		}

	case FieldDuration:
		var v float64
		if data == nil {
			val = NewDurationValueFromFloat64(v)
		} else if err = json.Unmarshal(*data, &v); err == nil {
			val = NewDurationValueFromFloat64(v)
		}

	case FieldCheckbox:
		var v bool
		if data == nil {
			val = NewCheckboxValue(v)
		} else if err = json.Unmarshal(*data, &v); err == nil {
			val = NewCheckboxValue(v)
		}

	case FieldAddress:
		var v string
		if data == nil {
			val = NewAddressValue(v)
		} else if err = json.Unmarshal(*data, &v); err == nil {
			val = NewAddressValue(v)
		}

	case FieldPhoneNumber:
		var v string
		if data == nil {
			val = NewPhoneNumberValue(v)
		} else if err = json.Unmarshal(*data, &v); err == nil {
			val = NewPhoneNumberValue(v)
		}

	case FieldEmailAddress:
		var v string
		if data == nil {
			val = NewEmailAddressValue(v)
		} else if err = json.Unmarshal(*data, &v); err == nil {
			val = NewEmailAddressValue(v)
		}

	case FieldUser:
		var v *User
		if data == nil {
			val = NewUserValue(v)
		} else if err = json.Unmarshal(*data, &v); err == nil {
			val = NewUserValue(v)
		}

	case FieldUserList:
		var v []*User
		if data == nil {
			val = NewListUserValue(v)
		} else if err = json.Unmarshal(*data, &v); err == nil {
			val = NewListUserValue(v)
		}

	case FieldFileAttachment:
		var v string
		if data == nil {
			val = NewEmailAddressValue(v)
		} else if err = json.Unmarshal(*data, &v); err == nil {
			val = NewEmailAddressValue(v)
		}

	case FieldReportLink:
		var v string
		if data == nil {
			val = NewReportLinkValue(v)
		} else if err = json.Unmarshal(*data, &v); err == nil {
			val = NewReportLinkValue(v)
		}

	case FieldURL:
		var v string
		if data == nil {
			val, err = NewURLValueFromString(v)
		} else if err = json.Unmarshal(*data, &v); err == nil {
			val, err = NewURLValueFromString(v)
		}

	case FieldiCalendar:
		var v string
		if data == nil {
			val = NewiCalendarValue(v)
		} else if err = json.Unmarshal(*data, &v); err == nil {
			val = NewiCalendarValue(v)
		}

	case FieldvCard:
		var v string
		if data == nil {
			val = NewPredecessorValue(v)
		} else if err = json.Unmarshal(*data, &v); err == nil {
			val = NewPredecessorValue(v)
		}

	case FieldPredecessor:
		var v string
		if data == nil {
			val = NewPredecessorValue(v)
		} else if err = json.Unmarshal(*data, &v); err == nil {
			val = NewPredecessorValue(v)
		}

	default:
		err = fmt.Errorf("unsupported field type (%s)", ftype)
	}

	if err != nil {
		err = fmt.Errorf("%s (fid %v)", err, fid)
	}

	return
}

// UnmarshalJSON implements json.UnmarshalJSON by using the field type to
// decode the "value" parameter into the appropraite data type.
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
	output.Data = make([]map[int]*QueryRecordsOutputData, len(v.Data))
	for i, record := range v.Data {

		data := make(map[int]*QueryRecordsOutputData, len(record))
		for fid, field := range record {

			// Get the Quick Base field type from the fid.
			ftype, ok := tmap[fid]
			if !ok {
				err = fmt.Errorf("field type not found (fid %v)", fid)
				return
			}

			// Unmarshal the field based on its Quick Base type.
			var val *Value
			if val, err = unmarshalField(fid, ftype, field.Value); err != nil {
				return
			}

			data[fid] = &QueryRecordsOutputData{Value: val}
		}

		output.Data[i] = data
	}

	return
}

type parseQueryRecordsOutput struct {
	ErrorProperties

	Fields   []*QueryRecordsOutputFields `json:"fields,omitempty"`
	Metadata *QueryRecordsOutputMetadata `json:"metadata,omitempty"`
	Data     []map[int]struct {
		Value *json.RawMessage `json:"value"`
	} `json:"data,omitempty"`
}

//
// Models for Quick Base field values.
//

// Timestamp models a unix timestamp in Quick Base.
type Timestamp struct {
	time.Time
}

// MarshalJSON converts time.Time to a unix timestamp in microseconds.
func (t Timestamp) MarshalJSON() ([]byte, error) {
	buf := bytes.NewBuffer([]byte{})
	json.NewEncoder(buf).Encode(t.Time.Format(FormatDateTime))
	return buf.Bytes(), nil
}

// UnmarshalJSON converts a unix timestamp in microseconds to a time.Time.
func (t *Timestamp) UnmarshalJSON(b []byte) error {

	var s string
	buf := bytes.NewBufferString(string(b))
	err := json.NewDecoder(buf).Decode(&s)
	if err != nil {
		return fmt.Errorf("error decoding timestamp (%s)", err)
	}

	t.Time, err = time.Parse(FormatDateTime, s)
	if err != nil {
		return fmt.Errorf("error parsing timestamp (%s)", err)
	}

	return nil
}

// User models a user in Quick Base.
type User struct {
	ID    string `json:"id"`
	Email string `json:"email"`
	Name  string `json:"name"`
}

//
// Field Properties
//

// Field models a field.
type Field struct {

	// Basics
	Label    string `json:"label,omitempty" validate:"required"`
	Type     string `json:"fieldType,omitempty" validate:"required"`
	Required bool   `json:"required,omitempty"`
	Unique   bool   `json:"unique,omitempty"`

	// Display
	DisplayInBold          bool `json:"bold,omitempty"`
	DisplayWithoutWrapping bool `json:"noWrap,omitempty"`

	// Advanced
	AutoFill        bool   `json:"doesDataCopy,omitempty"`
	Searchable      bool   `json:"findEnabled"`      // Defaults to true, so we cannot omitempty.
	AddToNewReports bool   `json:"appearsByDefault"` // Defaults to true, so we cannot omitempty.
	FieldHelpText   string `json:"fieldHelp,omitempty"`
	TrackField      bool   `json:"audited,omitempty"`

	// No UI
	AddToForms bool `json:"addToForms,omitempty"` // Not documented
}

// FieldProperties models field properties.
// TODO Make a custom unmarshaler to not show properties if the struct ie empty.
// SEE https://stackoverflow.com/a/28447372
type FieldProperties struct {

	// Basics
	DefaultValue string `json:"defaultValue,omitempty"`

	// Text - Multiple Choice field options
	AllowNewChoices    bool `json:"allowNewChoices,omitempty"`
	SortChoicesAsGiven bool `json:"sortAsGiven,omitempty"`

	// Display
	NumberOfLines   int `json:"numLines,omitempty"`
	MaxCharacters   int `json:"maxLength,omitempty"`
	WidthOfInputBox int `json:"width,omitempty"`

	// No UI
	ExactMatch   bool   `json:"exact,omitempty"`
	ForeignKey   bool   `json:"foreignKey,omitempty"`
	Formula      string `json:"formula,omitempty"`
	ParentTable  string `json:"masterTableTag,omitempty"`
	PrimaryKey   bool   `json:"primaryKey,omitempty"`
	RelatedField int    `json:"targetFieldId,omitempty"`

	// Comments
	Comments string `json:"comments,omitempty"`
}

// FieldPermission models the permissions properties.
type FieldPermission struct {
	Role   string `json:"role"`
	Type   string `json:"permissionType"`
	RoleID int    `json:"roleId"`
}

//
// Relationship Properties
//

// Relationship models a relationship.
type Relationship struct {
	ChildTableID    string               `json:"childTableId,omitempty"`
	ForeignKeyField *RelationshipField   `json:"foreignKeyField,omitempty"`
	RelationshipID  int                  `json:"id,omitempty"`
	IsCrossApp      bool                 `json:"isCrossApp,omitempty"`
	LookupFields    []*RelationshipField `json:"lookupFields,omitempty"`
	ParentTableID   string               `json:"parentTableId,omitempty"`
	SummaryFields   []*RelationshipField `json:"summaryFields,omitempty"`
}

// RelationshipField models fields in relationship output.
type RelationshipField struct {
	FieldID int    `json:"id,omitempty"`
	Label   string `json:"label,omitempty"`
	Type    string `json:"type,omitempty"`
}

// RelationshipSummaryField models summary fields in relationship input/output.
type RelationshipSummaryField struct {
	SummaryFieldID   int    `json:"summaryFid,omitempty"`
	Label            string `json:"label,omitempty"`
	AccumulationType string `json:"accumulationType,omitempty"`
	Where            string `json:"where,omitempty"`
}

func relationshipPath(tid string, rid int) string {
	return "/tables/" + url.PathEscape(tid) + "/relationship/" + strconv.Itoa(rid)
}

//
// App Properties
//

// App models an app.
// NOTE The description property is in ErrorProperties.
type App struct {
	AppID                    string      `json:"id,omitempty"`
	Name                     string      `json:"name,omitempty"`
	TimeZone                 string      `json:"timeZone,omitempty"`
	DateFormat               string      `json:"dateFormat,omitempty"`
	Created                  *Timestamp  `json:"created,omitempty"`
	Updated                  *Timestamp  `json:"updated,omitempty"`
	Variables                []*Variable `json:"variables,omitempty"`
	HasEveryoneOnTheInternet bool        `json:"hasEveryoneOnTheInternet,omitempty"`
}

// Variable models a variable.
type Variable struct {
	Name  string `xml:"-" json:"name,omitempty"`
	Value string `xml:"value,omitempty" json:"value,omitempty"`
}

//
// Input / Output interfaces and convenience functions.
//

// Input models the payload of API requests.
type Input interface {

	// url returns the URL the API request is sent to.
	url() string

	// method is the HTTP method used when sending API requests.
	method() string

	// addHeaders adds HTTP headers to the API request.
	addHeaders(req *http.Request)

	// encode encodes the request and writes it to io.Writer.
	encode() ([]byte, error)
}

// Output models the payload of API responses.
type Output interface {

	// decode parses the response in io.ReadCloser to Output.
	decode(io.ReadCloser) error

	// errorMessage returns the error message, if any.
	errorMessage() string

	// errorDetail returns the error detail, if any.
	errorDetail() string

	// handleError handles errors returned by the API.
	handleError(Output, *http.Response) error
}
