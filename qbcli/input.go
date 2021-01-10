package qbcli

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/QuickBase/quickbase-cli/qbclient"
)

// Prompt prompts a user for input and returns what they typed.
func Prompt(label string, validate qbclient.ValidateStringFn) (s string, err error) {
	reader := bufio.NewReader(os.Stdin)

	for {
		// Print the prompt label.
		fmt.Print(label)

		// Read the input.
		s, err = reader.ReadString('\n')
		if err != nil {
			return
		}

		// Trim spaces, inclusive of the newline.
		s = strings.TrimSpace(s)

		// Validate the input.
		if verr := validate(s); verr != nil {
			fmt.Printf("%s\n\n", verr)
		} else {
			break
		}
	}

	return
}
