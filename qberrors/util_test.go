package qberrors_test

import (
	"errors"
	"fmt"

	"github.com/QuickBase/quickbase-cli/qberrors"
)

func ExampleNotFoundError() {
	err := qberrors.NotFoundError("item %q", "123")

	if errors.Is(err, qberrors.NotFound) {
		fmt.Println(err)
	}

	// Output: item "123": not found
}
