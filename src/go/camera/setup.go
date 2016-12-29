package camera

import (
	"github.com/blackjack/webcam"
	"log"
	"image"
)

var cam *webcam.Webcam
var img image.Image
func SetupCam(c string) {			//Sets up camera with specified path
	log.Println("Opening camera "+c)
	kam,err := webcam.Open(c)
	panic(err)
	cam = kam
	var pfmt webcam.PixelFormat = 8	//Pixel format - V4L2_SRGB
	var h,w uint32 = 640,480	//Image size
	_,_,_,err = cam.SetImageFormat(pfmt,h,w)
	panic(err)
	panic(cam.StartStreaming())
	go frameLoop()
}
