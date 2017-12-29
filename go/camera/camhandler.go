package camera

import (
	"io"
	"net/http"
)

//CamHandler is an http.Handler for a camera
type CamHandler struct {
	Cam Camera
}

func (c CamHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	//dont cache the request
	w.Header().Set("Content-Type", "image/jpeg")
	w.Header().Set("Cache-Control", "max-age=0, no-cache, must-revalidate, proxy-revalidate")
	w.Header().Set("Pragma", "no-cache")
	w.Header().Set("Expires", "0")
	//get a frame
	fr, err := c.Cam.GetFrameJPEG()
	if err != nil {
		http.Error(w, "Error getting frame: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer fr.Close()
	//send frame
	_, err = io.Copy(w, fr)
	if err != nil {
		http.Error(w, "Relaying error "+err.Error(), http.StatusTeapot)
		return
	}
}
