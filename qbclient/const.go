package qbclient

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
