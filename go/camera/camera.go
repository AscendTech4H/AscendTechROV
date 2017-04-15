package camera

import (
	"archive/zip"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"

	"../startup"
)

var relayed [3]string

func multiCamHandler(w http.ResponseWriter, r *http.Request) {
	//Retrieve frames
	var frames [3][]byte
	for i := 0; i < 3; i++ {
		req, err := http.Get(fmt.Sprintf("http://localhost:8080/cam/%d", i))
		if err != nil {
			http.Error(w, fmt.Sprintf("Camera request error on camera %d: %s", i, err.Error()), http.StatusInternalServerError)
			return
		}
		defer req.Body.Close()
		resp, err := ioutil.ReadAll(req.Body)
		if err != nil {
			http.Error(w, fmt.Sprintf("Error reading request on camera %d: %s", i, err.Error()), http.StatusInternalServerError)
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
		flag.StringVar(&(relayed[0]), "cam0", "", "Camera connection 0")
		flag.StringVar(&(relayed[1]), "cam1", "", "Camera connection 1")
		flag.StringVar(&(relayed[2]), "cam2", "", "Camera connection 2")
		return nil
	})
	startup.NewTask(247, func() error {
		http.HandleFunc("/cam/", func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "image/jpeg")
			w.Header().Set("Cache-Control", "max-age=0, no-cache, must-revalidate, proxy-revalidate")
			w.Header().Set("Pragma", "no-cache")
			w.Header().Set("Expires", "0")

			strs := strings.Split(r.URL.Path, "/cam/")
			if len(strs) < 1 {
				http.Error(w, "Invalid camera URL", http.StatusBadRequest)
				return
			}
			camnum, err := strconv.Atoi(strs[len(strs)-1])
			if err != nil {
				http.Error(w, "Error decoding camera URL: "+err.Error(), http.StatusBadRequest)
				return
			}
			if (camnum > len(relayed)) || (camnum < 0) {
				http.Error(w, "Non-existant camera", http.StatusBadRequest)
				return
			}
			resp, err := http.Get(relayed[camnum])
			if err != nil {
				http.Error(w, "Relaying error "+err.Error(), http.StatusInternalServerError)
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
}
