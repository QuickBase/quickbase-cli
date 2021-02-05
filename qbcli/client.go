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

// FieldMap is a map of field IDs to field definitions.
type FieldMap map[int]*qbclient.ListFieldsOutputField

var _fmap map[string]FieldMap

// CacheTableSchema caches schema information for a table.
func CacheTableSchema(qb *qbclient.Client, tableID string) error {
	output, err := qb.ListFieldsByTableID(tableID)
	if err != nil {
		return err
	}

	m := make(FieldMap, len(output.Fields))
	for _, field := range output.Fields {
		m[field.FieldID] = field
	}

	_fmap[tableID] = m
	return nil
}

// GetTableSchema returns schema information for a table. If the schema is not
// in the in-memory cache, it retrieves the data and caches it.
func GetTableSchema(qb *qbclient.Client, tableID string) (FieldMap, error) {
	var err error
	_, ok := _fmap[tableID]
	if !ok {
		err = CacheTableSchema(qb, tableID)
	}
	return _fmap[tableID], err
}

// GetCachedTableSchema returns schema information for a table.
//
// TODO Caching beyond in-memory caching?
func GetCachedTableSchema(tableID string) (FieldMap, error) {
	m, ok := _fmap[tableID]
	if !ok {
		err := errors.New("field metadata not set")
		return FieldMap{}, fmt.Errorf("table %s: %w", tableID, err)
	}
	return m, nil
}

func init() {
	_fmap = make(map[string]FieldMap, 0)
}
