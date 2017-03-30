package camera

import (
	"flag"
	"log"
	"time"

	"../controller"
	"../copilot"
	"../debug"
	"../startup"
	"../util"

	"github.com/blackjack/webcam"
)

var camLoc string

func init() {
	startup.NewTask(1, func() error { //Set up can flag parsing
		flag.StringVar(&camLoc, "cam", "video0", "Cam connection")
		return nil
	})
	startup.NewTask(252, func() error { //Set up can flag parsing
		cam, err := webcam.Open(camLoc)
		util.UhOh(err)
		tick := time.NewTicker(time.Second / 2) //send the img twice a second
		go func() {
			for t := range tick.C {
				if debug.Verbose {
					log.Println("Start camera read at time " + t.String())
				}
				data, err := cam.ReadFrame()
				util.UhOh(err)
				controller.SendData(data)
				copilot.SendData(data)
			}
		}()
		return nil
	})
}
