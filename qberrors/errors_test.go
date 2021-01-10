package qberrors_test

import (
	"errors"
	"fmt"
	"net/http"
	"testing"

	"github.com/QuickBase/quickbase-cli/qberrors"
)

func TestErrors(t *testing.T) {
	up3 := errors.New("upstream")
	up4 := errors.New("unsafe")

	err1 := qberrors.SafeErrorf(qberrors.NotFound, "item %q", "123")
	err2 := fmt.Errorf("more context: %w", err1)
	err3 := fmt.Errorf("context: %w", qberrors.Internal(up3))
	err4 := up4
	err5 := qberrors.ErrSafe{Message: "test default status code"}

	tests := []struct {
		err          error
		exIsSafe     bool
		exStatusCode int
		exMessage    string
		exDetail     string
		exUpstream   error
	}{
		{err1, true, http.StatusNotFound, "not found", `item "123"`, nil},
		{err2, true, http.StatusNotFound, "not found", `more context: item "123"`, nil},
		{err3, true, http.StatusInternalServerError, "internal error", "context", up3},
		{err4, false, http.StatusInternalServerError, "internal error", "", up4},
		{err5, true, http.StatusInternalServerError, "test default status code", "", nil},
	}

	for _, tt := range tests {
		if actual := qberrors.IsSafe(tt.err); actual != tt.exIsSafe {
			t.Errorf("expected qberrors.IsSafe to return %t", tt.exIsSafe)
		}
		if actual := qberrors.StatusCode(tt.err); tt.exStatusCode != actual {
			t.Errorf("got %v, expected %v", actual, tt.exStatusCode)
		}
		if actual := qberrors.SafeMessage(tt.err); tt.exMessage != actual {
			t.Errorf("got %q, expected %q", actual, tt.exMessage)
		}
		if actual := qberrors.SafeDetail(tt.err); tt.exDetail != actual {
			t.Errorf("got %q, expected %q", actual, tt.exDetail)
		}
		if actual := qberrors.Upstream(tt.err); tt.exUpstream != actual {
			t.Errorf("got %q, expected %q", actual, tt.exUpstream)
		}
	}
}

func TestError(t *testing.T) {
	err := qberrors.Client(nil).Safe(qberrors.NotFound)

	if !errors.Is(err, qberrors.NotFound) {
		t.Error("expected errors.Is to be true")
	}

	if actual := err.Error(); actual != qberrors.NotFound.Message {
		t.Errorf("got %q, expected %q", actual, qberrors.NotFound.Message)
	}
}

func TestRetry(t *testing.T) {
	tests := []struct {
		err     qberrors.Error
		exRetry bool
	}{
		{qberrors.Client(nil), false},
		{qberrors.Internal(nil), false},
		{qberrors.Service(nil), true},
	}

	for _, tt := range tests {
		if actual := tt.err.Retry(); actual != tt.exRetry {
			t.Errorf("got %t, expected %t", actual, tt.exRetry)
		}
	}
}
