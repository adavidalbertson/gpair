package internal

import (
	"fmt"
)

// Verbose determines whether to print extra messages
var Verbose bool
// Help determines whether to print usage information
var Help bool

// PrintVerbose prints a formatted message if Verbose is true, otherwise it's a no-op
func PrintVerbose(format string, v ...interface{}) {
	if Verbose {
		fmt.Printf(format + "\n", v...)
	}
}