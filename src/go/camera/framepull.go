package camera

import (
	//"github.com/blackjack/webcam"
	//"log"
	"image"
)

func frameLoop() { //Loop that loads frames into img
	for {
		//Get frame
		panic(cam.WaitForFrame(1000))
		fr, err := cam.ReadFrame()
		panic(err)

		f := image.NewRGBA(image.Rect(0, 0, 640, 480))

		//Add alpha channel (RGB -> RGBA)
		p := 0
		pm := 640 * 480
		for p < pm {
			f.Pix[p*4] = fr[p*3]
			f.Pix[(p*4)+1] = fr[(p*3)+1]
			f.Pix[(p*4)+2] = fr[(p*3)+2]
			f.Pix[(p*4)+3] = 0
		}

		//Switch to new frame
		img = f
	}
}
