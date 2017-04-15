package main

import (
	"bytes"
	"flag"
	"fmt"
	"image"
	"image/jpeg"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"
	"sync"

	"github.com/blackjack/webcam"

	"../util"
)

var relayed [3]string

func init() {
	flag.StringVar(&(relayed[0]), "cam0", "null", "Camera connection 0")
	flag.StringVar(&(relayed[1]), "cam1", "null", "Camera connection 1")
	flag.StringVar(&(relayed[2]), "cam2", "null", "Camera connection 2")
	flag.IntVar(&quality, "jpg", 10, "JPEG quality")
	flag.Parse()
	Cams = make([]*Cam, len(relayed))
	for i, f := range relayed {
		if f != "null" {
			cam, err := OpenCam(f)
			util.UhOh(err)
			Cams[i] = cam
		}
	}
	http.HandleFunc("/cam/", camhandler)
}

var quality int

//Cam is a camera object
type Cam struct {
	cam *webcam.Webcam
	dat []byte
	lck sync.RWMutex
}

//FrameSend reads a frame from the camera
func (c *Cam) FrameSend(w io.Writer) {
	dat := func() []byte {
		c.lck.RLock()
		defer c.lck.RUnlock()
		return c.dat
	}()
	//Convert to jpeg
	img := image.NewYCbCr(image.Rect(0, 0, 640, 480), image.YCbCrSubsampleRatio422)
	pxgroups := make([]struct {
		Y1, Cb, Y2, Cr uint8
	}, (640*480)/2)
	for i := range pxgroups {
		v := &pxgroups[i]
		addr := i * 4
		slc := dat[addr : addr+4]
		v.Y1 = slc[0]
		v.Cb = slc[1]
		v.Y2 = slc[2]
		v.Cr = slc[3]
	}
	for i, v := range pxgroups {
		x := (i * 2) % 640
		x1, x2 := x, x+1
		y := (i * 2) / 640

		cpos := img.COffset(x, y)
		ypos1 := img.YOffset(x1, y)
		ypos2 := img.YOffset(x2, y)

		img.Cb[cpos] = v.Cb
		img.Cr[cpos] = v.Cr
		img.Y[ypos1] = v.Y1
		img.Y[ypos2] = v.Y2
	}
	util.UhOh(jpeg.Encode(w, img, &jpeg.Options{Quality: quality}))
}

//Frame gets a jpeg-encoded frame
func (c *Cam) Frame() []byte {
	buf := bytes.NewBuffer(nil)
	c.FrameSend(buf)
	return buf.Bytes()
}

//OpenCam opens a camera at a file path
func OpenCam(file string) (*Cam, error) {
	ca, err := webcam.Open(file)
	if err != nil {
		return nil, err
	}
	fmt.Println(ca.GetSupportedFormats())
	ca.SetImageFormat(webcam.PixelFormat(1448695129), 640, 480)
	err = ca.StartStreaming()
	if err != nil {
		return nil, err
	}
	cam := new(Cam)
	cam.cam = ca
	go func() { //Load frames in background
		for {
			err := cam.cam.WaitForFrame(1) //Up to one second to wait for a frame
			if err != nil {
				log.Fatalf("Camera %s crashed with error %s", file, err.Error())
			}
			dat, err := cam.cam.ReadFrame()
			if err != nil {
				cam.lck.Lock()
				cam.dat = nil
				cam.lck.Unlock()
				log.Fatalf("Camera %s crashed with error %s", file, err.Error())
			}
			cam.lck.Lock()
			cam.dat = dat
			cam.lck.Unlock()
		}
	}()
	return cam, nil
}

//Cams is an array of cameras used by the robot
var Cams []*Cam

func camhandler(writer http.ResponseWriter, requ *http.Request) {
	//Write as jpeg
	writer.Header().Set("Content-Type", "image/jpeg")
	writer.Header().Set("Cache-Control", "max-age=0, no-cache, must-revalidate, proxy-revalidate")
	writer.Header().Set("Pragma", "no-cache")
	writer.Header().Set("Expires", "0")
	strs := strings.Split(requ.URL.Path, "/cam/")
	if len(strs) < 1 {
		http.Error(writer, "Invalid camera URL", http.StatusBadRequest)
		return
	}
	camnum, err := strconv.Atoi(strs[len(strs)-1])
	if err != nil {
		http.Error(writer, "Error decoding camera URL: "+err.Error(), http.StatusBadRequest)
		return
	}
	if (camnum > len(Cams)) || (camnum < 0) {
		http.Error(writer, "Non-existant camera", http.StatusBadRequest)
		return
	}
	defer func() {
		if e := recover(); e != nil {
			http.Error(writer, err.Error(), http.StatusInternalServerError)
		}
	}()
	Cams[camnum].FrameSend(writer)
	writer.(http.Flusher).Flush()
}
func main() {
	util.UhOh(http.ListenAndServe(":8000", nil))
}
