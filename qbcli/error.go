package qbcli

import (
	"context"
	"os"

	"github.com/cpliakas/cliutil"
)

// HandleError handles an error by logging it and returning a non-zero status.
// We reserve Fatal errors for internal problems.
func HandleError(ctx context.Context, logger *cliutil.LeveledLogger, message string, err error) {
	if err != nil {
		logger.Error(ctx, message, err)
		os.Exit(1)
	}
}
