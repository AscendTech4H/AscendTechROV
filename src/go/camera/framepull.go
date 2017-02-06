package camera

import (
	//"github.com/blackjack/webcam"
	"log"
	"image"
	"../shmeh"
//	"fmt"
)

func frameLoop() { //Loop that loads frames into img
	for {
		//Get frame
		shmeh.UhOh(cam.WaitForFrame(1000))
		fr, err := cam.ReadFrame()
		shmeh.UhOh(err)

		f := image.NewRGBA(image.Rect(0, 0, 640, 480))
		//log.Println(fr)
		log.Println("out")
		log.Println(len(f.Pix)/4,len(fr)/3)

		//Add alpha channel (RGB -> RGBA)
		p := 0
		pm := 640 * 480
		for p < pm {
//			log.Println((p*4)+3,len(f.Pix))
//			log.Println((p*3)+2,len(fr))
			f.Pix[p*4] = fr[p*3]
			f.Pix[(p*4)+1] = fr[(p*3)+1]
			f.Pix[(p*4)+2] = fr[(p*3)+2]
			f.Pix[(p*4)+3] = 0
			p++
		}
		log.Println(pichans)
		//Switch to new frame
		pichans <- f
	}
}
