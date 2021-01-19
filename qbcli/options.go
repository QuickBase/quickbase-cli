package qbcli

import (
	"fmt"
	"reflect"
	"strconv"

	"github.com/QuickBase/quickbase-cli/qbclient"
	"github.com/cpliakas/cliutil"
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

func init() {
	cliutil.RegisterOptionTypeFunc("query", NewQueryOption)
	cliutil.RegisterOptionTypeFunc("record", NewRecordOption)
	cliutil.RegisterOptionTypeFunc("sort", NewSortOption)
	cliutil.RegisterOptionTypeFunc("group", NewGroupOption)
}
