package qbcli

import (
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"os"
	"sort"
	"time"

	"github.com/QuickBase/quickbase-cli/qbclient"
)

// ExportOptions are the options read through the command line.
type ExportOptions struct {
	TableID   string `validate:"required" cliutil:"option=table-id"`
	Filepath  string `cliutil:"option=file usage='file the data is exported to'"`
	BatchSize int    `cliutil:"option=batch-size default=10000"`
	Delay     int    `cliutil:"option=delay"`

	// Fields    []int  `cliutil:"option=fields"`
}

// Export exports data from a Quickbase table into an io.Writer.
func Export(qb *qbclient.Client, opts *ExportOptions) error {

	var file io.Writer
	if opts.Filepath != "" {
		var err error
		file, err = os.OpenFile(opts.Filepath, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644)
		if err != nil {
			return fmt.Errorf("error opening file: %w", err)
		}
	} else {
		file = os.Stdout
	}

	// Get the table's fields.
	fields, err := GetTableSchema(qb, opts.TableID)
	if err != nil {
		return fmt.Errorf("error getting table metadata: %w", err)
	}

	// Build a list of every fid.
	fids := []int{}
	for _, field := range fields {
		fids = append(fids, field.FieldID)
	}
	sort.Ints(fids)

	writer := csv.NewWriter(file)
	defer writer.Flush()

	// Write the header.
	header := make([]string, len(fids))
	for idx, fid := range fids {
		header[idx] = fields[fid].Label
	}
	writer.Write(header)

	// Batch read records.
	size := opts.BatchSize
	skip := 0
	for {

		// Query records, sorted by record ID.
		qri := &qbclient.QueryRecordsInput{
			Select: fids,
			From:   opts.TableID,
			SortBy: []*qbclient.QueryRecordsInputSortBy{
				{FieldID: 3, Order: qbclient.SortByASC},
			},
			Options: &qbclient.QueryRecordsInputOptions{
				Top:  size,
				Skip: skip,
			},
		}
		qro, err := qb.QueryRecords(qri)
		if err != nil {
			return fmt.Errorf("error querying records: %w", err)
		}

		// Write the row data.
		for _, record := range qro.Data {
			row := make([]string, len(fids))
			for idx, fid := range fids {
				row[idx] = record[fid].Value.String()
			}
			writer.Write(row)
		}

		// Flush the buffer.
		writer.Flush()
		if err = writer.Error(); err != nil {
			return err
		}

		// Increment skip for next page, break if on last page.
		skip += size
		if skip >= qro.Metadata.TotalRecords {
			break
		}

		// Delay before the next API call.
		if opts.Delay > 0 {
			time.Sleep(time.Duration(opts.Delay) * time.Millisecond)
		}
	}

	return nil
}

// ImportOptions are the options read through the command line.
type ImportOptions struct {
	TableID      string            `validate:"required" cliutil:"option=table-id"`
	Filepath     string            `cliutil:"option=file usage='file the data is imported from'"`
	BatchSize    int               `cliutil:"option=batch-size default=10000"`
	Map          map[string]string `cliutil:"option=map"`
	Delay        int               `cliutil:"option=delay"`
	Timeout      int               `cliutil:"option=timeout default=5 usage='timeout in seconds waiting for data to be read from stdin'"`
	MergeFieldID int               `cliutil:"option=merge-field-id"`

	// Fields    []int  `cliutil:"option=fields"`
}

// Import imports data from an io.Reader ito a Quickbase table.
func Import(qb *qbclient.Client, opts *ImportOptions) (*qbclient.InsertRecordsOutputMetadata, error) {
	metadata := &qbclient.InsertRecordsOutputMetadata{
		CreatedRecordIDs:              []int{},
		TotalNumberOfRecordsProcessed: 0,
		UnchangedRecordIDs:            []int{},
		UpdatedRecordIDs:              []int{},
	}

	var file io.Reader
	if opts.Filepath != "" {
		var err error
		file, err = os.Open(opts.Filepath)
		if err != nil {
			return metadata, fmt.Errorf("error opening file: %w", err)
		}
	} else {
		file = os.Stdin
		if err := waitStdin(opts.Timeout); err != nil {
			return metadata, err
		}
	}

	// Get the table's fields.
	fields, err := GetTableSchema(qb, opts.TableID)
	if err != nil {
		return metadata, fmt.Errorf("error getting table metadata: %w", err)
	}

	// Build a map of field label to fid.
	lmap := make(map[string]int, len(fields))
	for _, field := range fields {
		lmap[field.Label] = field.FieldID
	}

	reader := csv.NewReader(file)
	fmap := []int{}

	// Batch write records.
	line := 0
	eof := false
	records := []map[int]*qbclient.InsertRecordsInputData{}

	for {

		// Read each record from the CSV data.
		// TODO better error handling.
		row, err := reader.Read()
		if err == io.EOF {
			eof = true
		} else if err != nil {
			return metadata, fmt.Errorf("error reading line %v: %w", line, err)
		}

		// If first line, map the header to field IDs.
		// Otherwise, build the data records.
		if !eof {
			if line == 0 {
				for _, label := range row {

					// Check the field label map first.
					if destLabel, ok := opts.Map[label]; ok {
						label = destLabel
					}

					// Now get the field ID.
					fid, ok := lmap[label]
					if !ok {
						return metadata, fmt.Errorf("%s field not in destination table", label)
					}

					// Append the fid from the field map.
					fmap = append(fmap, fid)
				}
			} else {

				record := make(map[int]*qbclient.InsertRecordsInputData)
				for idx, data := range row {

					// TODO defensive coding ...
					fid := fmap[idx]

					// We cannot insert record metadata.
					if fid != 3 && fid <= 5 {
						continue
					}

					// Skip the record ID if it isn't the merge field.
					if fid == 3 && opts.MergeFieldID != 3 {
						continue
					}

					// TODO defensive coding ...
					ftype := fields[fid].Type

					// Create a *qbclient.Value from the string value and field type.
					val, err := qbclient.NewValueFromString(data, ftype)
					if err != nil {
						return metadata, fmt.Errorf("value invalid for field %v: %w", fid, err)
					}

					// Add the value to the record .
					record[fid] = &qbclient.InsertRecordsInputData{Value: val}
				}

				records = append(records, record)
			}
		}

		// Write the batch.
		if len(records) >= opts.BatchSize || (len(records) > 0 && eof) {
			input := &qbclient.InsertRecordsInput{
				To:           opts.TableID,
				Data:         records,
				MergeFieldID: opts.MergeFieldID,
			}

			iro, err := qb.InsertRecords(input)
			if err != nil {
				return metadata, fmt.Errorf("error inserting records: %w", err)
			}

			// Empty the records for the next batch.
			records = []map[int]*qbclient.InsertRecordsInputData{}

			metadata.CreatedRecordIDs = append(metadata.CreatedRecordIDs, iro.Metadata.CreatedRecordIDs...)
			metadata.TotalNumberOfRecordsProcessed += iro.Metadata.TotalNumberOfRecordsProcessed
			metadata.UnchangedRecordIDs = append(metadata.UnchangedRecordIDs, iro.Metadata.UnchangedRecordIDs...)
			metadata.UpdatedRecordIDs = append(metadata.UpdatedRecordIDs, iro.Metadata.UpdatedRecordIDs...)

			// Delay before the next API call.
			if opts.Delay > 0 && !eof {
				time.Sleep(time.Duration(opts.Delay) * time.Millisecond)
			}
		}

		// Break if we are at the end of the file.
		if eof {
			break
		}

		line++
	}

	return metadata, nil
}

// TODO move this to cliutil.
func waitStdin(wiat int) error {
	tick := time.Tick(100 * time.Millisecond)
	timeout := time.After(time.Duration(wiat) * time.Second)
	done := make(chan bool, 1)

	for {
		select {
		case <-tick:
			stat, err := os.Stdin.Stat()
			if err != nil {
				return fmt.Errorf("error getting info for stdin: %w", err)
			}
			if stat.Size() > 0 {
				done <- true
			}
		case <-done:
			return nil
		case <-timeout:
			return errors.New("timeout: no data read from stdin")
		}
	}
}
