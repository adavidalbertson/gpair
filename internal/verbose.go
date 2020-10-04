package internal

import (
	"fmt"
)

var Verbose bool

func PrintVerbose(format string, v ...interface{}) {
	if Verbose {
		fmt.Printf(format, v...)
	}
}