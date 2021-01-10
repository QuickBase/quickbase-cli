package qbclient_test

import (
	"testing"

	"github.com/QuickBase/quickbase-cli/qbclient"
)

func TestTableAlias(t *testing.T) {
	tables := []struct {
		label string
		want  string
	}{
		{"Table [WITH] some -chars", "_DBID_TABLE__WITH__SOME__CHARS"},
	}

	for _, tt := range tables {
		have := qbclient.TableAlias(tt.label)
		if have != tt.want {
			t.Errorf("have %q, want %q", have, tt.want)
		}
	}
}
