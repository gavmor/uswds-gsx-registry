// Command uswds-gsx vendors USWDS-flavored gsx components into a
// gsxui-initialized project by copying their source — you own the code.
// See https://github.com/gavmor/uswds-gsx-registry.
package main

import (
	"fmt"
	"os"

	"github.com/gavmor/uswds-gsx-registry/internal/cli"
)

func main() {
	if err := cli.Run(os.Args[1:]); err != nil {
		fmt.Fprintln(os.Stderr, "uswds-gsx:", err)
		os.Exit(1)
	}
}
