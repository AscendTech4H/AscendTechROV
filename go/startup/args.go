package startup

import (
	"flag"
)

func init() {
	NewTask(2, func() error { //Parse flags
		flag.Parse()
		return nil
	})
}
