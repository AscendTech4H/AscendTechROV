package failutil

import (
	"logger"
	"shutdown"
	"os"
)

//Function to crash
func CrashIfErr(log Logger,err error)
	if err!=nil {			//Error - clean up and exit
		log.Error(err)
		shutdown.Shutdown()
		os.Exit(1)
	}
}
