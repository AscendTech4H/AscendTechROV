package camera

import (
	"flag"
	"log"
	"net/http"
	"strconv"
	"strings"
	"sync"

	"../debug"
	"../startup"

	"github.com/blackjack/webcam"
)

var camLocs [3]string

//Cam is a camera object
type Cam struct {
	cam *webcam.Webcam
	dat []byte
	lck sync.RWMutex
}

//Frame reads a frame from the camera
func (c *Cam) Frame() []byte {
	c.lck.RLock()
	defer c.lck.RUnlock()
	return c.dat
}

//OpenCam opens a camera at a file path
func OpenCam(file string) (*Cam, error) {
	ca, err := webcam.Open(file)
	if err != nil {
		return nil, err
	}
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
	debug.VLog("Camera update request started")
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
	c := Cams[camnum]
	dat := c.Frame()
	if len(dat) < 1 {
		http.Error(writer, "Empty camera read", http.StatusInternalServerError)
		return
	}

	//We are good
	writer.WriteHeader(http.StatusOK)
	_, err = writer.Write(dat)
	if err != nil { //Not sure what would happen here
		debug.VLog("Write error: " + err.Error())
		debug.VLog("They call me teapot")
		http.Error(writer, "Write error: "+err.Error(), http.StatusTeapot)
		return
	}
	debug.VLog("Camera update request successfully completed")
}
func init() {
	startup.NewTask(1, func() error { //Set up can flag parsing
		flag.StringVar(&(camLocs[0]), "cam0", "/dev/video0", "Camera connection 0")
		flag.StringVar(&(camLocs[1]), "cam1", "/dev/video1", "Camera connection 1")
		flag.StringVar(&(camLocs[2]), "cam2", "/dev/video2", "Camera connection 2")
		return nil
	})
	startup.NewTask(245, func() error { //Open cameras
		Cams = make([]*Cam, len(camLocs))
		for i, f := range camLocs {
			if f != "null" {
				cam, err := OpenCam(f)
				if err != nil {
					return err
				}
				Cams[i] = cam
			}
		}
		return nil
	})
	startup.NewTask(247, func() error {
		http.HandleFunc("/cam/", camhandler)
		return nil
	})
}
