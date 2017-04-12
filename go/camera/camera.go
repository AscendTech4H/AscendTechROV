package camera

import (
	"archive/zip"
	"bytes"
	"flag"
	"fmt"
	"image"
	"image/jpeg"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"

	"../debug"
	"../startup"
	"../util"

	"github.com/blackjack/webcam"
)

var camLocs [2]string
var relayed string
var quality int

//Cam is a camera object
type Cam struct {
	cam *webcam.Webcam
	dat []byte
	lck sync.RWMutex
}

//FrameSend reads a frame from the camera
func (c *Cam) FrameSend(w io.Writer) {
	start := time.Now()
	debug.VLog("Waiting for image lock" + time.Since(start).String())
	start = time.Now()
	dat := func() []byte {
		c.lck.RLock()
		defer c.lck.RUnlock()
		return c.dat
	}()
	debug.VLog("Encoding image" + time.Since(start).String())
	start = time.Now()
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
	debug.VLog("Encoding jpeg" + time.Since(start).String())
	start = time.Now()
	util.UhOh(jpeg.Encode(w, img, &jpeg.Options{Quality: quality}))
	debug.VLog("done sending" + time.Since(start).String())
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
	defer func() {
		if e := recover(); e != nil {
			http.Error(writer, err.Error(), http.StatusInternalServerError)
		}
	}()
	Cams[camnum].FrameSend(writer)
	writer.(http.Flusher).Flush()
	debug.VLog("Camera update request successfully completed")
}
func multiCamHandler(w http.ResponseWriter, r *http.Request) {
	//Retrieve frames
	var frames [3][]byte
	for i := 0; i < 3; i++ {
		req, err := http.Get(fmt.Sprintf("localhost:8080/cam/%d", i))
		if err != nil {
			http.Error(w, fmt.Sprintf("Camera request error on camera %d: %s", i, err.Error()), http.StatusFailedDependency)
			return
		}
		defer req.Body.Close()
		resp, err := ioutil.ReadAll(req.Body)
		if err != nil {
			http.Error(w, fmt.Sprintf("Error reading request on camera %d: %s", i, err.Error()), http.StatusFailedDependency)
			return
		}
		frames[i] = resp
	}
	//Zip frames
	z := zip.NewWriter(w)
	for i := 0; i < 3; i++ {
		writer, err := z.Create(fmt.Sprintf("cam%d.jpg", i))
		if err != nil {
			http.Error(w, fmt.Sprintf("Zip file creation error: %s", err.Error()), http.StatusInternalServerError)
			return
		}
		_, err = writer.Write(frames[i])
		if err != nil {
			http.Error(w, fmt.Sprintf("Zip write error: %s", err.Error()), http.StatusInternalServerError)
			return
		}
	}
	err := z.Close()
	if err != nil {
		http.Error(w, fmt.Sprintf("Zip close error: %s", err.Error()), http.StatusInternalServerError)
		return
	}
}

func init() {
	startup.NewTask(1, func() error { //Set up can flag parsing
		flag.StringVar(&(camLocs[0]), "cam0", "/dev/video0", "Camera connection 0")
		flag.StringVar(&(camLocs[1]), "cam1", "/dev/video1", "Camera connection 1")
		flag.StringVar(&relayed, "cam2", "", "Relayed camera connection")
		flag.IntVar(&quality, "quality", 10, "Camera quality in percent (default: 10)")
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
		http.HandleFunc("/cam/2", func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "image/jpeg")
			w.Header().Set("Cache-Control", "max-age=0, no-cache, must-revalidate, proxy-revalidate")
			w.Header().Set("Pragma", "no-cache")
			w.Header().Set("Expires", "0")
			resp, err := http.Get(relayed)
			if err != nil {
				http.Error(w, "Relaying error "+err.Error(), http.StatusExpectationFailed)
				log.Println("Relaying error: " + err.Error())
				return
			}
			defer resp.Body.Close()
			_, err = io.Copy(w, resp.Body)
			if err != nil {
				http.Error(w, "Relaying error "+err.Error(), http.StatusTeapot)
				return
			}
		})
		http.HandleFunc("/cam/all", multiCamHandler)
		return nil
	})
	startup.NewTask(247, func() error {
		http.HandleFunc("/cam/", camhandler)
		return nil
	})
}
