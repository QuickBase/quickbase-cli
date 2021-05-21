package qbclient

import (
	"net/url"
	"strconv"
)

//
// Field Properties
//

// Field models a field.
type Field struct {
	// Create is true if we are creating the field.
	Create bool `json:"-"`

	// Basics
	Label    string `json:"label,omitempty" validate:"required_if=Create true" cliutil:"option=label"`
	Type     string `json:"fieldType,omitempty" validate:"required_if=Create true" cliutil:"option=type"`
	Required bool   `json:"required,omitempty" cliutil:"option=required"`
	Unique   bool   `json:"unique,omitempty" cliutil:"option=unique"`

	// Display
	DisplayInBold          bool `json:"bold,omitempty" cliutil:"option=bold"`
	DisplayWithoutWrapping bool `json:"noWrap,omitempty" cliutil:"option=no-wrap"`

	// Advanced
	AutoFill        bool   `json:"doesDataCopy,omitempty" cliutil:"option=auto-fill"`
	Searchable      bool   `json:"findEnabled" cliutil:"option=searchable default=true func=boolstring"`
	AddToNewReports bool   `json:"appearsByDefault" cliutil:"option=add-to-reports default=true func=boolstring"`
	FieldHelpText   string `json:"fieldHelp,omitempty" cliutil:"option=help-text"`
	TrackField      bool   `json:"audited,omitempty" cliutil:"option=track-field"`

	// No UI
	AddToForms bool `json:"addToForms,omitempty" cliutil:"option=add-to-forms"`
}

// FieldProperties models field properties.
// TODO Make a custom unmarshaler to not show properties if the struct ie empty.
// SEE https://stackoverflow.com/a/28447372
type FieldProperties struct {

	// Basics
	DefaultValue string `json:"defaultValue,omitempty" cliutil:"option=default"`

	// Text - Multiple Choice field options
	AllowNewChoices    bool `json:"allowNewChoices,omitempty" cliutil:"option=allow-new-choices"`
	SortChoicesAsGiven bool `json:"sortAsGiven,omitempty" cliutil:"option=sort-as-given"`

	// Display
	NumberOfLines   int `json:"numLines,omitempty" cliutil:"option=num-lines"`
	MaxCharacters   int `json:"maxLength,omitempty" cliutil:"option=max"`
	WidthOfInputBox int `json:"width,omitempty" cliutil:"option=width"`

	// No UI
	ExactMatch   bool   `json:"exact,omitempty" cliutil:"option=exact-match"`
	ForeignKey   bool   `json:"foreignKey,omitempty" cliutil:"option=foreign-key"`
	Formula      string `json:"formula,omitempty" cliutil:"option=formula"`
	FormulaFile  string `json:"-" cliutil:"option=formula-file func=ioreader"`
	ParentTable  string `json:"masterTableTag,omitempty" cliutil:"option=parent-table"`
	PrimaryKey   bool   `json:"primaryKey,omitempty" cliutil:"option=primary-key"`
	RelatedField int    `json:"targetFieldId,omitempty" cliutil:"option=related-field"`

	// Comments
	Comments string `json:"comments,omitempty" cliutil:"option=comments"`
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
	SummaryFieldID   int    `json:"summaryFid,omitempty" cliutil:"option=field-id"`
	Label            string `json:"label,omitempty" cliutil:"option=label"`
	AccumulationType string `json:"accumulationType,omitempty" cliutil:"option=accumulation-type"`
	Where            string `json:"where,omitempty" cliutil:"option=where"`
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
// Report properties
//

// Report models the report object.
type Report struct {
	ReportID   string            `json:"id,omitempty"`
	Name       string            `json:"name,omitempty"`
	Type       string            `json:"type,omitempty"`
	Query      *ReportQuery      `json:"query,omitempty"`
	Properties *ReportProperties `json:"properties,omitempty"`
	UsedLast   *Timestamp        `json:"usedLast,omitempty"`
	UsedCount  int               `json:"usedCount,omitempty"`
	OwnerID    int               `json:"ownerId,omitempty"`
}

// ReportWithDescripton is a Report with a Description property. This prevents
// conflicts with ErrorProperties.
type ReportWithDescripton struct {
	Report
	Description string `json:"description,omitempty"`
}

// ReportProperties models the properties object.
type ReportProperties struct {
	DisplayOnlyNewOrChanged   bool                          `json:"displayOnlyNewOrChangedRecords,omitempty"`
	Column                    *ReportPropertiesColumn       `json:"columnProperties,omitempty"`
	StartDate                 *Date                         `json:"startDate,omitempty"`
	EndDate                   *Date                         `json:"endDate,omitempty"`
	StartingFieldID           int                           `json:"startingFieldId,omitempty"`
	EndingFieldID             int                           `json:"endingFieldId,omitempty"`
	MilestoneFieldID          int                           `json:"milestoneFieldId,omitempty"`
	SortByStartingField       bool                          `json:"sortByStartingField,omitempty"`
	Crosstab                  *ReportPropertiesCrosstab     `json:"crosstabs,omitempty"`
	Summmarize                *[]ReportPropertiesSummmarize `json:"summmarize,omitempty"`
	SortBy                    *[]ReportPropertiesSortBy     `json:"sortBy,omitempty"`
	AddressFieldID            int                           `json:"addressFieldId,omitempty"`
	MapViewType               string                        `json:"mapViewType,omitempty"`
	Stages                    []string                      `json:"stages,omitempty"`
	ManualOrdering            bool                          `json:"manualOrdering,omitempty"`
	StagesFromFieldID         int                           `json:"stagesFromFieldId,omitempty"`
	StartDateFieldID          int                           `json:"startDateField,omitempty"`
	EndDateFieldID            int                           `json:"endDateFieldId,omitempty"`
	DataLabel                 string                        `json:"dataLabel,omitempty"`
	Categories                *ReportPropertiesCategories   `json:"categories,omitempty"`
	DataSources               []*ReportPropertiesDataSource `json:"dataSources,omitempty"`
	ChartType                 string                        `json:"chartType,omitempty"`
	Series                    *ReportPropertiesSeries       `json:"series,omitempty"`
	Bubbles                   *ReportPropertiesBubbles      `json:"bubbles,omitempty"`
	DataSourcesLabel          string                        `json:"dataSourcesLabel,omitempty"`
	SecondaryDataSourcesLabel string                        `json:"secondaryDataSourcesLabel,omitempty"`
	SecondaryDataSources      []*ReportPropertiesDataSource `json:"secondaryDataSources,omitempty"`
	LineLabel                 string                        `json:"lineLabel,omitempty"`
	BarLabel                  string                        `json:"barLabel,omitempty"`
	BarDataSources            []*ReportPropertiesDataSource `json:"barDataSources,omitempty"`
	LineDataSources           []*ReportPropertiesDataSource `json:"lineDataSources,omitempty"`
	Goal                      *ReportPropertiesGoal         `json:"goal,omitempty"`
	Range                     *ReportPropertiesRange        `json:"range,omitempty"`
}

// ReportPropertiesColumn models the columnProperties object.
type ReportPropertiesColumn struct {
	FieldID       int    `json:"fieldId,omitempty"`
	LabelOverride string `json:"labelOverride,omitempty"`
}

// ReportPropertiesCrosstab models the crosstabs object.
type ReportPropertiesCrosstab struct {
	FieldID  int    `json:"fieldId,omitempty"`
	Grouping string `json:"grouping,omitempty"`
}

// ReportPropertiesSummmarize models the summmarize object.
type ReportPropertiesSummmarize struct {
	Type        string `json:"type,omitempty"`
	FieldID     int    `json:"fieldId,omitempty"`
	Aggregation string `json:"aggregation,omitempty"`
	ShowAs      int    `json:"showAs,omitempty"`
}

// ReportPropertiesSortBy models the sortBy object.
type ReportPropertiesSortBy struct {
	ElementIndex int    `json:"summarizationElementIndex,omitempty"`
	Order        string `json:"order,omitempty"`
	By           string `json:"by,omitempty"`
}

// ReportPropertiesCategories models the categories object.
type ReportPropertiesCategories struct {
	FieldID  int    `json:"fieldId,omitempty"`
	Label    string `json:"label,omitempty"`
	Grouping string `json:"grouping,omitempty"`
}

// ReportPropertiesDataSource models the dataSource object.
type ReportPropertiesDataSource struct {
	Type        string `json:"type,omitempty"`
	FieldID     int    `json:"fieldId,omitempty"`
	Label       string `json:"label,omitempty"`
	Aggregation string `json:"aggregation,omitempty"`
}

// ReportPropertiesSeries models the series object.
type ReportPropertiesSeries struct {
	FieldID  int    `json:"fieldId,omitempty"`
	Grouping string `json:"grouping,omitempty"`
}

// ReportPropertiesBubbles models the bubbles object.
type ReportPropertiesBubbles struct {
	Type        string `json:"type,omitempty"`
	FieldID     int    `json:"fieldId,omitempty"`
	Aggregation string `json:"aggregation,omitempty"`
}

// ReportPropertiesGoal models the goal object.
type ReportPropertiesGoal struct {
	Value int    `json:"number,omitempty"`
	Label string `json:"label,omitempty"`
}

// ReportPropertiesRange models the range object.
type ReportPropertiesRange struct {
	Type                string  `json:"type,omitempty"`
	SmallestValue       float64 `json:"smallestValue,omitempty"`
	LargestValue        float64 `json:"largestValue,omitempty"`
	FieldID             int     `json:"fieldId,omitempty"`
	CalculationMethod   string  `json:"calculationMethod,omitempty"`
	IntervalsPercentage []int   `json:"intervalsPercentage,omitempty"`
}

// ReportQuery models the query object.
type ReportQuery struct {
	TableID       string                      `json:"tableId,omitempty"`
	Filter        string                      `json:"filter,omitempty"`
	FormulaFields []*ReportQueryFormulaFields `json:"formulaFields,omitempty"`
	Fields        []int                       `json:"fields,omitempty"`
	SortBy        []*ReportQuerySortBy        `json:"sortBy,omitempty"`
	GroupBy       []*ReportQueryGroupBy       `json:"groupBy,omitempty"`
}

// ReportQueryFormulaFields models the formulaFields
type ReportQueryFormulaFields struct {
	FormulaFieldID   int    `json:"id,omitempty"`
	Label            string `json:"label,omitempty"`
	Type             string `json:"fieldType,omitempty"`
	Formula          string `json:"formula,omitempty"`
	DecimalPrecision int    `json:"decimalPrecision,omitempty"`
}

// ReportQuerySortBy models the sortBy object.
type ReportQuerySortBy struct {
	FieldID int    `json:"fieldId,omitempty"`
	Order   string `json:"order,omitempty"`
}

// ReportQueryGroupBy models the sortBy object.
type ReportQueryGroupBy struct {
	FieldID  int    `json:"fieldId,omitempty"`
	Grouping string `json:"grouping,omitempty"`
}
