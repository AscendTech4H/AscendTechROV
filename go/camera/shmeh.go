package camera

import (
	//	"github.com/blackjack/webcam"
	"image"
	//	"log"
)

var pichans chan image.Image
var pichanr chan chan image.Image
var fchan chan func()

func picloop() {
	sc <- true
	o := []func(){}
	var pic image.Image
	for {
		select {
		case pic = <-pichans:
		case f := <-fchan:
			o = append(o, f)
		case c := <-pichanr:
			c <- pic
			for _, f := range o {
				go f()
			}
		}
	}
}

// Get current frame
func GetPic() image.Image {
	c := make(chan image.Image)
	pichanr <- c
	return <-c
}

// Add update handler
func OnFrameUpdate(f func()) {
	fchan <- f
}
