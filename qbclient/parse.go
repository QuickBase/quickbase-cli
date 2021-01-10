package qbclient

import (
	"bytes"
	"encoding/csv"
	"io"
)

// ParseList parses a string into a string slice.
func ParseList(in string) (out []string, err error) {
	buf := bytes.NewBufferString(in)
	r := csv.NewReader(buf)
	out, err = r.Read()
	if err == io.EOF {
		err = nil
	}
	return
}
