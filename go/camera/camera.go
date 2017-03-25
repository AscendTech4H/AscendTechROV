package controller

import (
	"../controller"
	"../copilot"
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
		go func() {
			for {
				sleep(500) //send the img twice a second
				data, err := cam.ReadFrame()
				util.UhOh(err)
				controller.SendData(data)
				copilot.SendData(data)
			}
		}()
	})
}
