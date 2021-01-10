package qbcli

import (
	"github.com/QuickBase/quickbase-cli/qbclient"
	"github.com/cpliakas/cliutil"
	"github.com/spf13/viper"
)

// FieldOption* constants contain common field option names.
const (
	FieldOptionAddToForms             = "add-to-forms"
	FieldOptionAddToNewReports        = "add-to-new-reports"
	FieldOptionAllowNewChoices        = "allow-new-choices"
	FieldOptionAutoFill               = "auto-fill"
	FieldOptionDefaultValue           = "default-value"
	FieldOptionDisplayInBold          = "bold"
	FieldOptionDisplayWithoutWrapping = "no-wrap"
	FieldOptionExactMatch             = "exact-match"
	FieldOptionFieldHelpText          = "help-text"
	FieldOptionForeignKey             = "foreign-key"
	FieldOptionFormula                = "formula"
	FieldOptionLabel                  = "label"
	FieldOptionMaxCharacters          = "max-chars"
	FieldOptionNumberOfLines          = "num-lines"
	FieldOptionRequired               = "required"
	FieldOptionParentTable            = "parent-table"
	FieldOptionPrimaryKey             = "primary-key"
	FieldOptionRelatedField           = "related-field"
	FieldOptionSearchable             = "searchable"
	FieldOptionSortChoicesAsGiven     = "sort-as-given"
	FieldOptionTrackField             = "track"
	FieldOptionType                   = "type"
	FieldOptionUnique                 = "unique"
	FieldOptionWidthOfInputBox        = "width"
)

// SetFieldOptions sets common field options.
func SetFieldOptions(flags *cliutil.Flagger) {

	// Basics
	flags.String(FieldOptionLabel, "", "", "")
	flags.String(FieldOptionType, "", "", "")
	flags.Bool(FieldOptionRequired, "", false, "")
	flags.Bool(FieldOptionUnique, "", false, "")

	// Display
	flags.Bool(FieldOptionDisplayInBold, "", false, "")
	flags.Bool(FieldOptionDisplayWithoutWrapping, "", false, "")

	// Advanced
	flags.Bool(FieldOptionAutoFill, "", false, "")
	flags.Bool(FieldOptionSearchable, "", false, "")
	flags.Bool(FieldOptionAddToNewReports, "", false, "")
	flags.String(FieldOptionFieldHelpText, "", "", "")
	flags.Bool(FieldOptionTrackField, "", false, "")

	// No UI
	flags.Bool(FieldOptionAddToForms, "", false, "")

	// Properties

	// Basics
	flags.String(FieldOptionDefaultValue, "", "", "")

	// Text - Multiple Choice field options
	flags.Bool(FieldOptionAllowNewChoices, "", false, "")
	flags.Bool(FieldOptionSortChoicesAsGiven, "", false, "")

	// Display
	flags.Int(FieldOptionNumberOfLines, "", 0, "")
	flags.Int(FieldOptionMaxCharacters, "", 0, "")
	flags.Int(FieldOptionWidthOfInputBox, "", 0, "")

	// No UI
	flags.Bool(FieldOptionExactMatch, "", false, "")
	flags.Bool(FieldOptionForeignKey, "", false, "")
	flags.String(FieldOptionFormula, "", "", "")
	flags.String(FieldOptionParentTable, "", "", "")
	flags.Bool(FieldOptionPrimaryKey, "", false, "")
	flags.Int(FieldOptionRelatedField, "", 0, "")
}

// NewFieldFromOptions returns a new qbclient.Field from values in cfg.
func NewFieldFromOptions(cfg *viper.Viper) (f qbclient.Field) {

	// Basics
	f.Type = cfg.GetString(FieldOptionType)
	f.Label = cfg.GetString(FieldOptionLabel)
	cliutil.SetBoolValue(cfg, FieldOptionRequired, &f.Required)
	cliutil.SetBoolValue(cfg, FieldOptionUnique, &f.Unique)

	// Display
	cliutil.SetBoolValue(cfg, FieldOptionDisplayInBold, &f.DisplayInBold)
	cliutil.SetBoolValue(cfg, FieldOptionDisplayWithoutWrapping, &f.DisplayWithoutWrapping)

	// Advanced
	cliutil.SetBoolValue(cfg, FieldOptionAutoFill, &f.AutoFill)
	cliutil.SetBoolValue(cfg, FieldOptionSearchable, &f.Searchable)
	cliutil.SetBoolValue(cfg, FieldOptionAddToNewReports, &f.AddToNewReports)
	cliutil.SetStringValue(cfg, FieldOptionFieldHelpText, &f.FieldHelpText)
	cliutil.SetBoolValue(cfg, FieldOptionTrackField, &f.TrackField)

	return
}

// NewPropertiesFromOptions returns a new qbclient.Properties from values in cfg.
// See https://developer.quickbase.com/operation/createField
func NewPropertiesFromOptions(cfg *viper.Viper) (p qbclient.FieldProperties) {

	// Basics
	cliutil.SetStringValue(cfg, FieldOptionDefaultValue, &p.DefaultValue)

	// Text - Multiple Choice field options
	cliutil.SetBoolValue(cfg, FieldOptionAllowNewChoices, &p.AllowNewChoices)
	cliutil.SetBoolValue(cfg, FieldOptionSortChoicesAsGiven, &p.SortChoicesAsGiven)

	// Display
	cliutil.SetIntValue(cfg, FieldOptionNumberOfLines, &p.NumberOfLines)
	cliutil.SetIntValue(cfg, FieldOptionMaxCharacters, &p.MaxCharacters)
	cliutil.SetIntValue(cfg, FieldOptionWidthOfInputBox, &p.WidthOfInputBox)

	// No UI
	cliutil.SetBoolValue(cfg, FieldOptionExactMatch, &p.ExactMatch)
	cliutil.SetBoolValue(cfg, FieldOptionForeignKey, &p.ForeignKey)
	cliutil.SetStringValue(cfg, FieldOptionFormula, &p.Formula)
	cliutil.SetStringValue(cfg, FieldOptionParentTable, &p.ParentTable)
	cliutil.SetBoolValue(cfg, FieldOptionPrimaryKey, &p.PrimaryKey)
	cliutil.SetIntValue(cfg, FieldOptionRelatedField, &p.RelatedField)

	return
}
