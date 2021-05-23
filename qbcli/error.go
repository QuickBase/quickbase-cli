package qbcli

import (
	"context"
	"net/http"
	"os"

	"github.com/QuickBase/quickbase-cli/qberrors"
	"github.com/cpliakas/cliutil"
)

var (
	TestsFailed = qberrors.ErrSafe{Message: "tests failed", StatusCode: http.StatusBadRequest}
)

func TestsFailedError(format string, a ...interface{}) error {
	return qberrors.Client(nil).Safef(TestsFailed, format, a...)
}

// HandleError handles an error by logging it and returning a non-zero status.
// We reserve Fatal errors for internal problems.
func HandleError(ctx context.Context, logger *cliutil.LeveledLogger, message string, err error) {
	if err != nil {
		logger.Error(ctx, message, err)
		os.Exit(1)
	}
}
