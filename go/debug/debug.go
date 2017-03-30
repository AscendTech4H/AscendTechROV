package debug

import (
	"flag"

	"../startup"
)

//Verbose indicates whether verbose debugging is enabled
var Verbose bool

func init() {
	startup.NewTask(1, func() error { //Set up can flag parsing
		flag.BoolVar(&Verbose, "verbose", false, "enable verbose debugging")
		return nil
	})
}
