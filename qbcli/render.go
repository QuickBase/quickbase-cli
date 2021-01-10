package qbcli

import (
	"context"
	"errors"
	"fmt"

	"github.com/QuickBase/quickbase-cli/qberrors"
	"github.com/cpliakas/cliutil"
	"github.com/spf13/cobra"
)

// Render renders the output in JSON, or writes an error log.
func Render(
	ctx context.Context,
	logger *cliutil.LeveledLogger,
	cmd *cobra.Command,
	cfg GlobalConfig,
	v interface{},
	err error,
) {

	// Render the error.
	if err != nil {
		ctx = cliutil.ContextWithLogTag(ctx, "code", fmt.Sprintf("%v", qberrors.StatusCode(err)))
		HandleError(ctx, logger, qberrors.SafeMessage(err), errors.New(qberrors.SafeDetail(err)))
	}

	// Render the output.
	if !cfg.Quiet() {
		err := cliutil.PrintJSONWithFilter(v, cfg.JMESPathFilter())
		HandleError(ctx, logger, "invalid JMESPath filter", err)
	}
}
