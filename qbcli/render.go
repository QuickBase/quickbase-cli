package qbcli

import (
	"context"
	"errors"
	"fmt"
	"reflect"

	"github.com/QuickBase/quickbase-cli/qbclient"
	"github.com/QuickBase/quickbase-cli/qberrors"
	"github.com/cpliakas/cliutil"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/spf13/cobra"
)

// Render renders the output in JSON, or writes an error log.
func Render(
	ctx context.Context,
	logger *cliutil.LeveledLogger,
	cmd *cobra.Command,
	cfg GlobalConfig,
	v interface{},
	err error,
) {

	// Render the error.
	if err != nil {
		ctx = cliutil.ContextWithLogTag(ctx, "code", fmt.Sprintf("%v", qberrors.StatusCode(err)))
		HandleError(ctx, logger, qberrors.SafeMessage(err), errors.New(qberrors.SafeDetail(err)))
	}

	// Do not return output.
	if cfg.Quiet() {
		return
	}

	// Try to render a table.
	if cfg.Format() == "table" || cfg.Format() == "csv" || cfg.Format() == "markdown" {
		rerr := renderTable(v, cfg.Format())
		HandleError(ctx, logger, "error rendering table", rerr)
		return
	}

	// Default to rendering JSON.
	rerr := cliutil.PrintJSONWithFilter(v, cfg.JMESPathFilter())
	HandleError(ctx, logger, "JMESPath filter not valid", rerr)
}

func renderTable(a interface{}, format string) error {
	tw := table.NewWriter()

	// Only pointers!
	// This will panic otherwise. This is an internal function, but we whould
	// be a little more defensive to prevent that from happening.
	rv := reflect.ValueOf(a).Elem()
	rt := reflect.TypeOf(a).Elem()

	for idx := 0; idx < rt.NumField(); idx++ {
		rvf := rv.Field(idx)
		if !rvf.CanInterface() {
			continue
		}

		// Look for the embedded Report field.
		i := rvf.Interface()
		switch i.(type) {
		case qbclient.Records:
			r := i.(qbclient.Records)

			// map of field ids to index position in the table.
			fmap := make(map[int]int, len(r.Fields))

			// Add the header.
			header := make(table.Row, len(r.Fields))
			for idx, f := range r.Fields {
				header[idx] = f.Label
				fmap[f.FieldID] = idx
			}
			tw.AppendHeader(header)

			// Add the table data.
			data := make([]table.Row, len(r.Data))
			for idx, row := range r.Data {
				data[idx] = make(table.Row, len(row))
				for fid, record := range row {
					data[idx][fmap[fid]] = record.Value.String()
				}
			}
			tw.AppendRows(data)
		}
	}

	switch format {
	case "table":
		fmt.Println(tw.Render())
	case "csv":
		fmt.Println(tw.RenderCSV())
	case "markdown":
		fmt.Println(tw.RenderMarkdown())
	default:
		return fmt.Errorf("%s: format not valid", format)
	}

	return nil
}
