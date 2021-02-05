package qbclient

import (
	"bytes"
	"encoding/csv"
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"strconv"
	"time"

	"github.com/araddon/dateparse"
)

// ErrInvalidType is an invalid field type error.
var ErrInvalidType = errors.New("field type invalid")

// Value models the value of fields in Quick Base. This struct effectively
// handles the Quickbase field type / Golang type transformations.
type Value struct {
	Bool      bool
	Duration  time.Duration
	File      *File
	Float64   float64
	Str       string
	StrSlice  []string
	Time      time.Time
	URL       *url.URL
	User      *User
	UserSlice []*User

	QuickBaseType string
}

// String returns Value as a string.
func (v *Value) String() string {
	switch v.QuickBaseType {

	case FieldRecordID:
		return strconv.FormatFloat(v.Float64, 'f', -1, 64)

	case FieldText, FieldTextMultiLine, FieldTextMultipleChoice, FieldRichText:
		// Also picks up:
		// FieldAddressStreet1, FieldAddressStreet2, FieldAddressCity,
		// FieldAddressStateRegion, FieldAddressPostalCode, and
		// FieldAddressCountry
		return v.Str

	case FieldMultiSelectText:
		return writeCSV(v.StrSlice)

	case FieldNumeric, FieldNumericCurrency, FieldNumericPercent, FieldNumericRating:
		return strconv.FormatFloat(v.Float64, 'f', -1, 64)

	case FieldDate:
		return v.Time.UTC().Format(FormatDate)

	case FieldDateTime:
		return v.Time.UTC().Format(FormatDateTime)

	case FieldTimeOfDay:
		return v.Time.UTC().Format(FormatTimeOfDay)

	case FieldDuration:
		return fmt.Sprint(v.Duration.Milliseconds())

	case FieldCheckbox:
		return fmt.Sprintf("%t", v.Bool)

	case FieldAddress:
		return v.Str

	case FieldPhoneNumber:
		return v.Str

	case FieldEmailAddress:
		return v.Str

	case FieldUser:
		return v.User.ID

	case FieldUserList:
		s := make([]string, len(v.UserSlice))
		for idx, u := range v.UserSlice {
			s[idx] = u.ID
		}
		return writeCSV(s)

	case FieldFileAttachment:
		return v.File.URL

	case FieldReportLink:
		return v.Str

	case FieldURL:
		return v.URL.String()

	// Deprecated
	case FieldiCalendar:
		return v.Str

	// Deprecated
	case FieldvCard:
		return v.Str

	case FieldPredecessor:
		return v.Str

	default:
		return ""
	}
}

func writeCSV(s []string) string {
	buf := bytes.NewBuffer([]byte(``))
	w := csv.NewWriter(buf)
	w.Write(s)
	w.Flush()
	return buf.String()
}

// Timestamp models a Quickbase timestamp.
type Timestamp struct {
	time.Time
}

// MarshalJSON formats time.Time according to FormatDateTime.
func (t Timestamp) MarshalJSON() ([]byte, error) {
	return marshalTime(t.Time, FormatDateTime)
}

// UnmarshalJSON converts a timestamp to a time.Time.
func (t *Timestamp) UnmarshalJSON(b []byte) (err error) {
	t.Time, err = unmarshalTime(b, FormatDateTime)
	return
}

// Date models a Quickbase date.
type Date struct {
	time.Time
}

// MarshalJSON formats time.Time according to FormatDate.
func (d Date) MarshalJSON() ([]byte, error) {
	return marshalTime(d.Time, FormatDate)
}

// UnmarshalJSON converts a date to a time.Time.
func (d *Date) UnmarshalJSON(b []byte) (err error) {
	d.Time, err = unmarshalTime(b, FormatDate)
	return
}

func marshalTime(tm time.Time, layout string) ([]byte, error) {
	buf := bytes.NewBuffer([]byte{})
	json.NewEncoder(buf).Encode(tm.Format(layout))
	return buf.Bytes(), nil
}

func unmarshalTime(b []byte, layout string) (tm time.Time, err error) {
	var s string
	buf := bytes.NewBufferString(string(b))
	err = json.NewDecoder(buf).Decode(&s)
	if err != nil {
		err = fmt.Errorf("error decoding timestamp: %w", err)
		return
	}

	tm, err = time.Parse(layout, s)
	if err != nil {
		err = fmt.Errorf("error parsing timestamp: %w", err)
		return
	}

	return
}

// User models a user in Quick Base.
type User struct {
	ID    string `json:"id"`
	Email string `json:"email"`
	Name  string `json:"name"`
}

// File models a file attachment.
type File struct {
	URL     string         `json:"url,omitempty"`
	Version []*FileVersion `json:"versions,omitempty"`
}

// FileVersion models the "version" property.
type FileVersion struct {
	Creator  *FileCreator `json:"creator"`
	FileName string       `json:"name"`
	Uploaded *Timestamp   `json:"uploaded"`
	Version  int          `json:"versionNumber"`
}

// FileCreator models the "version.creator" property.
type FileCreator struct {
	Email  string `json:"email"`
	Name   string `json:"name"`
	UserID string `json:"id"`
}

// NewRecordIDValue returns a new Value of the FieldRecordID type.
func NewRecordIDValue(val float64) *Value {
	return &Value{Float64: val, QuickBaseType: FieldRecordID}
}

// NewRecordIDValueFromString returns a new Value of the FieldRecordID type.
func NewRecordIDValueFromString(val string) (*Value, error) {
	return parseStringToNumericValue(val, FieldRecordID)
}

// NewTextValue returns a new Value of the FieldText type.
func NewTextValue(val string) *Value {
	return &Value{Str: val, QuickBaseType: FieldText}
}

// NewTextMultiLineValue returns a new Value of the FieldTextMultiLine type.
func NewTextMultiLineValue(val string) *Value {
	return &Value{Str: val, QuickBaseType: FieldTextMultiLine}
}

// NewTextMultipleChoiceValue returns a new Value of the FieldTextMultipleChoice type.
func NewTextMultipleChoiceValue(val string) *Value {
	return &Value{Str: val, QuickBaseType: FieldTextMultipleChoice}
}

// NewRichTextValue returns a new Value of the FieldRichText type.
func NewRichTextValue(val string) *Value {
	return &Value{Str: val, QuickBaseType: FieldRichText}
}

// NewMultiSelectTextValue returns a new Value of the FieldMultiSelectText type.
func NewMultiSelectTextValue(val []string) *Value {
	return &Value{StrSlice: val, QuickBaseType: FieldMultiSelectText}
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
	return &Value{Str: val, QuickBaseType: FieldAddress}
}

// NewAddressStreet1Value returns a new Value of the FieldAddressStreet1 type.
func NewAddressStreet1Value(val string) *Value {
	return &Value{Str: val, QuickBaseType: FieldAddressStreet1}
}

// NewAddressStreet2Value returns a new Value of the FieldAddressStreet2 type.
func NewAddressStreet2Value(val string) *Value {
	return &Value{Str: val, QuickBaseType: FieldAddressStreet2}
}

// NewAddressCityValue returns a new Value of the FieldAddressCity type.
func NewAddressCityValue(val string) *Value {
	return &Value{Str: val, QuickBaseType: FieldAddressCity}
}

// NewAddressStateRegionValue returns a new Value of the FieldAddressStateRegion type.
func NewAddressStateRegionValue(val string) *Value {
	return &Value{Str: val, QuickBaseType: FieldAddressStateRegion}
}

// NewAddressPostalCodeValue returns a new Value of the FieldAddressPostalCode type.
func NewAddressPostalCodeValue(val string) *Value {
	return &Value{Str: val, QuickBaseType: FieldAddressPostalCode}
}

// NewAddressCountryValue returns a new Value of the FieldAddressCountry type.
func NewAddressCountryValue(val string) *Value {
	return &Value{Str: val, QuickBaseType: FieldAddressCountry}
}

// NewPhoneNumberValue returns a new Value of the FieldPhoneNumber type.
func NewPhoneNumberValue(val string) *Value {
	return &Value{Str: val, QuickBaseType: FieldPhoneNumber}
}

// NewEmailAddressValue returns a new Value of the FieldEmailAddress type.
func NewEmailAddressValue(val string) *Value {
	return &Value{Str: val, QuickBaseType: FieldEmailAddress}
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
func NewFileAttachmentValue(val *File) *Value {
	return &Value{File: val, QuickBaseType: FieldFileAttachment}
}

// NewReportLinkValue returns a new Value of the FieldReportLink type.
func NewReportLinkValue(val string) *Value {
	return &Value{Str: val, QuickBaseType: FieldReportLink}
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
	return &Value{Str: val, QuickBaseType: FieldiCalendar}
}

// NewvCardValue returns a new Value of the FieldvCard type.
// Make this a URL?
func NewvCardValue(val string) *Value {
	return &Value{Str: val, QuickBaseType: FieldvCard}
}

// NewPredecessorValue returns a new Value of the FieldPredecessor type.
func NewPredecessorValue(val string) *Value {
	return &Value{Str: val, QuickBaseType: FieldPredecessor}
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

	case FieldRecordID:
		v, err = NewRecordIDValueFromString(val)

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
		err = errors.New("cannot create file attachment field from string")

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
		err = fmt.Errorf("%s: %w", ftype, ErrInvalidType)
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
		return json.Marshal(v.Str)

	case FieldMultiSelectText:
		return json.Marshal(v.StrSlice)

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
		return json.Marshal(v.Str)

	case FieldPhoneNumber:
		return json.Marshal(v.Str)

	case FieldEmailAddress:
		return json.Marshal(v.Str)

	case FieldUser:
		return json.Marshal(v.User)

	case FieldUserList:
		return json.Marshal(v.UserSlice)

	case FieldFileAttachment:
		return json.Marshal(v.File)

	case FieldReportLink:
		return json.Marshal(v.Str)

	case FieldURL:
		return json.Marshal(v.URL.String())

	// Deprecated
	case FieldiCalendar:
		return json.Marshal(v.Str)

	// Deprecated
	case FieldvCard:
		return json.Marshal(v.Str)

	case FieldPredecessor:
		return json.Marshal(v.Str)

	default:
		return []byte(``), fmt.Errorf("%s: %w", v.QuickBaseType, ErrInvalidType)
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
		var v *File
		if data == nil {
			val = NewFileAttachmentValue(v)
		} else if err = json.Unmarshal(*data, &v); err == nil {
			val = NewFileAttachmentValue(v)
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
		err = fmt.Errorf("%s: %w", ftype, ErrInvalidType)
	}

	if err != nil {
		err = fmt.Errorf("fid %v: %w", fid, err)
	}

	return
}
