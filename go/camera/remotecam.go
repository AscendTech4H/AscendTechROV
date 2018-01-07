package camera

import (
	"errors"
	"io"
	"net/http"
	"net/url"
	"sync"
)

//Remote is a remote camera
type Remote struct {
	lck  sync.RWMutex
	targ string
}

//GetFrameJPEG returns an io.ReadCloser with a frame in JPEG format
func (r *Remote) GetFrameJPEG() (io.ReadCloser, error) {
	r.lck.RLock()
	defer r.lck.RUnlock()
	if r.targ == "" {
		return nil, errors.New("camera closed")
	}
	g, err := http.Get(r.targ)
	if err != nil {
		return nil, err
	}
	return g.Body, nil
}

//Close a camera.Remote
func (r *Remote) Close() error {
	r.lck.Lock()
	defer r.lck.Unlock()
	if r.targ == "" {
		return errors.New("already closed")
	}
	r.targ = ""
	return nil
}

//NewRemoteCam creates a new Remote object
func NewRemoteCam(u *url.URL) *Remote {
	return &Remote{targ: u.String()}
}
