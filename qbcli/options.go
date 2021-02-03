package qbcli

import (
	"context"
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"strings"

	"github.com/QuickBase/quickbase-cli/qbclient"
	"github.com/cpliakas/cliutil"
	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	en_translations "github.com/go-playground/validator/v10/translations/en"
	"github.com/spf13/viper"
)

// QueryOption implements Option for string options that contain queries.
type QueryOption struct {
	tag map[string]string
}

// NewQueryOption is a cliutil.OptionTypeFunc that returns a *cliutil.QueryOption.
func NewQueryOption(tag map[string]string) cliutil.OptionType { return &QueryOption{tag} }

// Set implements cliutil.OptionType.Set.
func (opt *QueryOption) Set(f *cliutil.Flagger) error {
	f.String(opt.tag["option"], opt.tag["short"], opt.tag["default"], opt.tag["usage"])
	return nil
}

// Read implements cliutil.OptionType.Read.
func (opt *QueryOption) Read(cfg *viper.Viper, field reflect.Value) error {
	s := ParseQuery(cfg.GetString(opt.tag["option"]))
	field.SetString(s)
	return nil
}

// SortOption implements Option for string options that contain queries.
type SortOption struct {
	tag map[string]string
}

// NewSortOption is a cliutil.OptionTypeFunc that returns a *cliutil.SortOption.
func NewSortOption(tag map[string]string) cliutil.OptionType { return &SortOption{tag} }

// Set implements cliutil.OptionType.Set.
func (opt *SortOption) Set(f *cliutil.Flagger) error {
	f.String(opt.tag["option"], opt.tag["short"], opt.tag["default"], opt.tag["usage"])
	return nil
}

// Read implements cliutil.OptionType.Read.
func (opt *SortOption) Read(cfg *viper.Viper, field reflect.Value) error {
	s := cfg.GetString(opt.tag["option"])
	if s == "" {
		return nil
	}

	v, err := ParseSortBy(s)
	if err != nil {
		return err
	}

	field.Set(reflect.ValueOf(v))
	return nil
}

// GroupOption implements Option for string options that contain queries.
type GroupOption struct {
	tag map[string]string
}

// NewGroupOption is a cliutil.OptionTypeFunc that returns a *cliutil.GroupOption.
func NewGroupOption(tag map[string]string) cliutil.OptionType { return &GroupOption{tag} }

// Set implements cliutil.OptionType.Set.
func (opt *GroupOption) Set(f *cliutil.Flagger) error {
	f.String(opt.tag["option"], opt.tag["short"], opt.tag["default"], opt.tag["usage"])
	return nil
}

// Read implements cliutil.OptionType.Read.
func (opt *GroupOption) Read(cfg *viper.Viper, field reflect.Value) error {
	s := cfg.GetString(opt.tag["option"])
	if s == "" {
		return nil
	}

	v, err := ParseGroupBy(s)
	if err != nil {
		return err
	}

	field.Set(reflect.ValueOf(v))
	return nil
}

// RecordOption implements Option for string options that contain record datas.
type RecordOption struct {
	tag map[string]string
}

// NewRecordOption is a cliutil.OptionTypeFunc that returns a *cliutil.RecordOption.
func NewRecordOption(tag map[string]string) cliutil.OptionType { return &RecordOption{tag} }

// Set implements cliutil.OptionType.Set.
func (opt *RecordOption) Set(f *cliutil.Flagger) error {
	f.String(opt.tag["option"], opt.tag["short"], opt.tag["default"], opt.tag["usage"])
	return nil
}

// Read implements cliutil.OptionType.Read.
func (opt *RecordOption) Read(cfg *viper.Viper, field reflect.Value) error {

	// Get the table's field type map.
	// TODO Do we need to remove the hard-coding to "to"?
	m, err := GetFieldTypeMap(cfg.GetString("to"))
	if err != nil {
		return err
	}

	// Parse the key/value pairs into *qbclient.Value objects, and build the
	// record being inserted into the table.
	record := make(map[int]*qbclient.InsertRecordsInputData)
	data := cliutil.ParseKeyValue(cfg.GetString(opt.tag["option"]))
	for k, v := range data {

		// Make sure k is an integer.
		fid, err := strconv.Atoi(k)
		if err != nil {
			return fmt.Errorf("key %s: %w", k, err)
		}

		// Get the field type from the table's field type map.
		ft, ok := m[fid]
		if !ok {
			return fmt.Errorf("field %v not defined in table", fid)
		}

		// Create a *qbclient.Value from the string value and field type.
		val, err := qbclient.NewValueFromString(v, ft)
		if err != nil {
			return fmt.Errorf("value invalid for field %v: %w", fid, err)
		}

		// Add the value to the record .
		record[fid] = &qbclient.InsertRecordsInputData{Value: val}
	}

	// Set the field's value.
	r := []map[int]*qbclient.InsertRecordsInputData{record}
	field.Set(reflect.ValueOf(r))
	return nil
}

// GetOptions gets options based on the input and validates them.
func GetOptions(ctx context.Context, logger *cliutil.LeveledLogger, input interface{}, cfg *viper.Viper) {
	err := cliutil.GetOptions(input, cfg)
	logger.FatalIfError(ctx, "error getting options", err)

	validate := validator.New()
	english := en.New()
	uni := ut.New(english, english)
	trans, _ := uni.GetTranslator("en")
	_ = en_translations.RegisterDefaultTranslations(validate, trans)

	// Custom translation for the "required" validator.
	validate.RegisterTranslation("required", trans, func(ut ut.Translator) error {
		return ut.Add("required", "{0} option is required", true)
	}, func(ut ut.Translator, fe validator.FieldError) string {
		// TODO We should be defensive, even if the error conditions shouldn't happen.
		field, _ := reflect.ValueOf(input).Elem().Type().FieldByName(fe.Field())
		tag := cliutil.ParseKeyValue(field.Tag.Get("cliutil"))
		t, _ := ut.T("required", tag["option"])
		return t
	})

	// Other validators we need to translate:
	//
	// - required_if (See Field.Label)
	// - min (See DeleteFieldsInput.FieldID)

	msgs := []string{}
	verr := validate.Struct(input)
	if verr != nil {
		verrs := verr.(validator.ValidationErrors)
		for _, ve := range verrs {
			msgs = append(msgs, ve.Translate(trans))
		}
	}

	if len(msgs) > 0 {
		HandleError(ctx, logger, "input not valid", errors.New(strings.Join(msgs, ", ")))
	}
}

func init() {
	cliutil.RegisterOptionTypeFunc("query", NewQueryOption)
	cliutil.RegisterOptionTypeFunc("record", NewRecordOption)
	cliutil.RegisterOptionTypeFunc("sort", NewSortOption)
	cliutil.RegisterOptionTypeFunc("group", NewGroupOption)

	cliutil.SetOptionMetadata("app-id", map[string]string{"usage": "the app's unique identifier, e.g., bqgruir3g"})
	cliutil.SetOptionMetadata("child-table-id", map[string]string{"usage": "the child table's unique identifier, e.g., bqgruir7z"})
	cliutil.SetOptionMetadata("data", map[string]string{"usage": "the record data in key=value format, e.g., '6=\"Another Record\" 7=3'"})
	cliutil.SetOptionMetadata("field-id", map[string]string{"usage": "the fields's unique identifier, e.g., 6"})
	cliutil.SetOptionMetadata("fields-to-return", map[string]string{"usage": "the list/range of fields to return, e.g., 6,7,10:15"})
	cliutil.SetOptionMetadata("from", map[string]string{"usage": "the table's unique identifier, e.g., bqgruir7z"})
	cliutil.SetOptionMetadata("group-by", map[string]string{"usage": "group records by fields, e.g., '6 DESC,7 ASC,8 equal-values'"})
	cliutil.SetOptionMetadata("lookup-field-ids", map[string]string{"usage": "the list/range of fids for lookup fields to create, e.g., 6,7,10:15"})
	cliutil.SetOptionMetadata("parent-table-id", map[string]string{"usage": "the parent table's unique identifier, e.g., bqgruir6f"})
	cliutil.SetOptionMetadata("relationship-id", map[string]string{"usage": "the relationship's unique identifier, e.g., 10"})
	cliutil.SetOptionMetadata("report-id", map[string]string{"usage": "the report's unique identifier, e.g., 1"})
	cliutil.SetOptionMetadata("select", map[string]string{"usage": "the list/range of fields to return, e.g., 6,7,10:15"})
	cliutil.SetOptionMetadata("skip", map[string]string{"usage": "the number of records to skip"})
	cliutil.SetOptionMetadata("sort-by", map[string]string{"usage": "sort records by fields, e.g., '6 DESC,8 ASC'"})
	cliutil.SetOptionMetadata("table-id", map[string]string{"usage": "the table's unique identifier, e.g., bqgruir7z"})
	cliutil.SetOptionMetadata("to", map[string]string{"usage": "the table's unique identifier, e.g., bqgruir7z"})
	cliutil.SetOptionMetadata("top", map[string]string{"usage": "the maximum number of records to display"})
	cliutil.SetOptionMetadata("use-app-time", map[string]string{"usage": "run the query against a date time field with the app's local time instead of UTC"})
	cliutil.SetOptionMetadata("where", map[string]string{"usage": "the filter, using the Quickbase query language or simplified syntax e.g., {3.EX.2}, 3=2"})
}
