package camera

import (
	"github.com/blackjack/webcam"
	"image"
	"log"
)

var cam *webcam.Webcam
var img image.Image
var sc chan bool

func SetupCam(c string) {
	pichans = make(chan image.Image)
        pichanr = make(chan chan image.Image)
        fchan = make(chan func())
	sc = make(chan bool)
	log.Println("Opening camera " + c)
	kam, err := webcam.Open(c)
	if err!=nil {
		panic(err)
	}
	cam = kam
	var pfmt webcam.PixelFormat = 8 //Pixel format - V4L2_SRGB
	var h, w uint32 = 640, 480      //Image size
	_, _, _, err = cam.SetImageFormat(pfmt, h, w)
	log.Println(cam.GetSupportedFormats())
	if err!=nil {
		panic(err)
	}
	if err = cam.StartStreaming(); err!=nil {
		panic(err)
	}
	go picloop()
	<-sc
	go frameLoop()
}
