package camera

import (
	"io"
	"net/http"
	"net/url"
)

//Remote is a remote camera
type Remote struct {
	targ string
}

func (r *Remote) GetFrameJPEG() (io.ReadCloser, error) {
	g, err := http.Get(r.targ)
	if err != nil {
		return nil, err
	}
	return g.Body, nil
}

func (r *Remote) Close() error {
	return nil
}

//NewRemoteCam creates a new Remote object
func NewRemoteCam(u *url.URL) *Remote {
	return &Remote{targ: u.String()}
}
