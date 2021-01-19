package qbcli

import (
	"context"
	"errors"
	"fmt"
	"os"

	"github.com/QuickBase/quickbase-cli/qbclient"
	"github.com/cpliakas/cliutil"
	"github.com/rs/xid"
	"github.com/spf13/cobra"
)

// NewLogger returns a new *cliutil.LeveledLogger.
func NewLogger(cmd *cobra.Command, cfg GlobalConfig) (ctx context.Context, logger *cliutil.LeveledLogger, transid xid.ID) {
	ctx, logger, transid = cliutil.NewLoggerWithContext(context.Background(), cfg.LogLevel())
	logger.SetOutput(os.Stderr)

	// Open the log file and set the logger to write to it.
	if logFile := cfg.LogFile(); logFile != "" {
		file, err := os.OpenFile(logFile, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)
		HandleError(ctx, logger, "error opening log file", err)
		logger.SetOutput(file)
	}

	return
}

// NewClient returns a new *qbclient.Client.
func NewClient(cmd *cobra.Command, cfg GlobalConfig) (ctx context.Context, logger *cliutil.LeveledLogger, qb *qbclient.Client) {
	var transid xid.ID
	ctx, logger, transid = NewLogger(cmd, cfg)

	// Instantiate the Quick Base API client with the logger plugin.
	qb = qbclient.New(cfg)
	qb.AddPlugin(NewLoggerPlugin(ctx, logger))

	// Dump raw requests and responses to the dump directory.
	if dumpDir := cfg.DumpDirectory(); dumpDir != "" {
		qb.AddPlugin(NewDumpPlugin(ctx, logger, transid.String(), dumpDir))
	}

	return
}

var tmap map[string]map[int]string

// SetFieldTypeMap sets the field type map for a table.
func SetFieldTypeMap(qb *qbclient.Client, tableID string) error {
	output, err := qb.ListFieldsByTableID(tableID)
	if err != nil {
		return err
	}

	m := make(map[int]string, len(output.Fields))
	for _, field := range output.Fields {
		m[field.FieldID] = field.Type
	}

	tmap[tableID] = m
	return nil
}

// GetFieldTypeMap returns a mapping of field IDs to Quickbase field type for
// all the fields in a table.
//
// TODO Caching?
func GetFieldTypeMap(tableID string) (map[int]string, error) {
	m, ok := tmap[tableID]
	if !ok {
		err := errors.New("field type map not set")
		return map[int]string{}, fmt.Errorf("table %s: %w", tableID, err)
	}
	return m, nil
}

func init() {
	tmap = make(map[string]map[int]string, 0)
}
