package qbclient

import (
	"fmt"
	"strings"
)

// Field* constants contain the Quick Base field types.
const (
	FieldRecordID           = "recordid"
	FieldText               = "text"
	FieldTextMultiLine      = "text-multi-line"
	FieldTextMultipleChoice = "text-multiple-choice"
	FieldRichText           = "rich-text"
	FieldMultiSelectText    = "multitext"
	FieldNumeric            = "numeric"
	FieldNumericCurrency    = "currency"
	FieldNumericPercent     = "percent"
	FieldNumericRating      = "rating"
	FieldDate               = "date"
	FieldDateTime           = "timestamp"
	FieldTimeOfDay          = "timeofday"
	FieldDuration           = "duration"
	FieldCheckbox           = "checkbox"
	FieldAddress            = "address"
	FieldAddressStreet1     = "text"
	FieldAddressStreet2     = "text"
	FieldAddressCity        = "text"
	FieldAddressStateRegion = "text"
	FieldAddressPostalCode  = "text"
	FieldAddressCountry     = "text"
	FieldPhoneNumber        = "phone"
	FieldEmailAddress       = "email"
	FieldUser               = "userid"
	FieldUserList           = "multiuserid"
	FieldFileAttachment     = "file"
	FieldURL                = "url"
	FieldReportLink         = "dblink"
	FieldiCalendar          = "ICalendarButton"
	FieldvCard              = "vCardButton"
	FieldPredecessor        = "predecessor"
)

// FieldType returns the Quick Base field type from in, which is a string that
// contains the constant without the "Field" prefix.
func FieldType(in string) (out string, err error) {
	in = strings.ToLower(strings.ReplaceAll(in, "_", ""))

	switch in {
	case "recordid":
		out = FieldRecordID
	case "text":
		out = FieldText
	case "textmultiline":
		out = FieldTextMultiLine
	case "textmultiplechoice":
		out = FieldTextMultipleChoice
	case "richtext":
		out = FieldRichText
	case "multiselecttext", FieldMultiSelectText:
		out = FieldMultiSelectText
	case FieldNumeric:
		out = FieldNumeric
	case "numericcurrency", FieldNumericCurrency:
		out = FieldNumericCurrency
	case "numericpercent", FieldNumericPercent:
		out = FieldNumericPercent
	case "numericrating", FieldNumericRating:
		out = FieldNumericRating
	case "date":
		out = FieldDate
	case "datetime", FieldDateTime:
		out = FieldDateTime
	case "timeofday":
		out = FieldTimeOfDay
	case "duration":
		out = FieldDuration
	case "checkbox":
		out = FieldCheckbox
	case "address":
		out = FieldAddress
	case "addressstreet1":
		out = FieldAddressStreet1
	case "addressstreet2":
		out = FieldAddressStreet2
	case "addresscity":
		out = FieldAddressCity
	case "addressstateregion":
		out = FieldAddressStateRegion
	case "addresspostalcode":
		out = FieldAddressPostalCode
	case "addresscountry":
		out = FieldAddressCountry
	case "phonenumber", FieldPhoneNumber:
		out = FieldPhoneNumber
	case "emailaddress", FieldEmailAddress:
		out = FieldEmailAddress
	case "user", FieldUser:
		out = FieldUser
	case "userlist", FieldUserList:
		out = FieldUserList
	case "fileattachment", FieldFileAttachment:
		out = FieldFileAttachment
	case "url":
		out = FieldURL
	case "reportlink", FieldReportLink:
		out = FieldReportLink
	case "icalendar", FieldiCalendar:
		out = FieldiCalendar
	case "vcard", FieldvCard:
		out = FieldvCard
	case "predecessor":
		out = FieldPredecessor
	default:
		err = fmt.Errorf("type not valid (%s)", in)
	}

	return
}

// AccumulationType* constants contain valid accumulation types for summary
// fields.
const (
	AccumulationTypeAverage           = "AVG"
	AccumulationTypeSum               = "SUM"
	AccumulationTypeMaximum           = "MAX"
	AccumulationTypeMinimum           = "MIN"
	AccumulationTypeStandardDeviation = "STD-DEV"
	AccumulationTypeCount             = "COUNT"
	AccumulationTypeCombinedText      = "COMBINED-TEXT"
	AccumulationTypeDistinctCount     = "DISTINCT-COUNT"
)

// Format* constants contain common format strings.
const (
	FormatDate      = "2006-01-02"
	FormatDateTime  = "2006-01-02T15:04:05Z"
	FormatTimeOfDay = "15:04:05"
)

// SortBy* constants model values used in the the order property.
const (
	SortByASC  = "ASC"
	SortByDESC = "DESC"
)
