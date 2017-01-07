package camera

import (
	"github.com/blackjack/webcam"
	"image"
	"log"
)

var pichans chan image.Image
var pichanr chan chan image.Image

func picloop() {
	pichans = make(chan image.Image)
	pichanr = make(chan chan image.Image)
	var pic image.Image
	for {
		select {
		case pic = <-pichans:
		case c := <-pichanr:
			c <- pic
		}
	}
}

func GetPic() image.Image {
	c := make(chan image.Image)
	pichanr <- c
	return <-c
}
