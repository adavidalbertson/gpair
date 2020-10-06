package internal

import (
	"fmt"
)

var Verbose bool
var Help bool

func PrintVerbose(format string, v ...interface{}) {
	if Verbose {
		fmt.Printf(format, v...)
	}
}