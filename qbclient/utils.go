package qbclient

import (
	"regexp"
	"strings"
)

// TableAlias converts a label to a table alias.
func TableAlias(label string) string {
	re := regexp.MustCompile(`\W`)
	return strings.ToUpper("_DBID_" + re.ReplaceAllString(label, "_"))
}
