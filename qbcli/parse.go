package qbcli

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/QuickBase/quickbase-cli/qbclient"
	"github.com/cpliakas/cliutil"
)

var reSortBy, reGroupBy *regexp.Regexp

// ParseFieldList parses a list of integers from a string.
func ParseFieldList(s string) (fids []int, err error) {
	if s == "" {
		return
	}

	parts := strings.Split(s, ",")
	fids = make([]int, len(parts))

	for i, part := range parts {
		part = strings.TrimSpace(part)
		if !cliutil.IsNumber(part) {
			err = errors.New("expecting fid to be a number")
			return
		}
		fid, _ := strconv.Atoi(part)
		fids[i] = fid
	}

	return
}

// ParseQuery parses queries. It also detcts and transforms simple queries into
// Quick Base query syntax.
func ParseQuery(q string) string {

	// Returns as-is if using Quick Base query syntax.
	if strings.HasPrefix(q, "{") || strings.HasPrefix(q, "(") {
		return q
	}

	// Parse the key/value pairs.
	m := cliutil.ParseKeyValue(q)
	clauses := make([]string, len(m))

	// Convert simple syntax into Quick Base query syntax.
	i := 0
	for k, v := range m {
		if v == "" {
			clauses[i] = fmt.Sprintf("{3.EX.%q}", k)
		} else {
			clauses[i] = fmt.Sprintf("{%q.EX.%q}", k, v)
		}
		i++
	}

	// Join all clauses by AND.
	return strings.Join(clauses, " AND ")
}

// ParseSortBy parses the sortBy clause.
func ParseSortBy(s string) (sortBy []*qbclient.QueryRecordsInputSortBy, err error) {
	clauses := strings.Split(s, ",")
	sortBy = make([]*qbclient.QueryRecordsInputSortBy, len(clauses))

	for i, clause := range clauses {
		matches := reSortBy.FindAllStringSubmatch(clause, -1)

		// Return an error if any clause cannot be parsed.
		if len(matches) == 0 {
			err = errors.New("no match")
			return
		}

		// Convert the field ID to an integer. Panic on any errors, because the
		// regex should only parse integers. Signifies logic error in app.
		fid, err := strconv.Atoi(matches[0][1])
		if err != nil {
			panic(err)
		}

		// Default to ASC.
		order := matches[0][2]
		if order == "" {
			order = qbclient.SortByASC
		}

		sortBy[i] = &qbclient.QueryRecordsInputSortBy{
			FieldID: fid,
			Order:   order,
		}
	}

	return
}

// ParseGroupBy parses the groupBy clause.
func ParseGroupBy(s string) (groupBy []*qbclient.QueryRecordsInputGroupBy, err error) {
	clauses := strings.Split(s, ",")
	groupBy = make([]*qbclient.QueryRecordsInputGroupBy, len(clauses))

	for i, clause := range clauses {
		matches := reGroupBy.FindAllStringSubmatch(clause, -1)

		// Return an error if any clause cannot be parsed.
		if len(matches) == 0 {
			err = errors.New("no match")
			return
		}

		// Convert the field ID to an integer. Panic on any errors, because the
		// regex should only parse integers. Signifies logic error in app.
		fid, err := strconv.Atoi(matches[0][1])
		if err != nil {
			panic(err)
		}

		groupBy[i] = &qbclient.QueryRecordsInputGroupBy{
			FieldID:  fid,
			Grouping: matches[0][2],
		}
	}

	return
}

func init() {
	reSortBy = regexp.MustCompile(`^\s*(\d+)(?:\s+(ASC|DESC)?\s*)?$`)
	reGroupBy = regexp.MustCompile(`^\s*(\d+)(?:\s+([-A-Za-z0-9_]+)?\s*)?$`)
}
