package camera

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"sync"

	"github.com/blackjack/webcam"
)

//Local is a local camera
type Local struct {
	curframe []byte
	toclose  bool
	closed   bool
	lck      sync.RWMutex
	done     sync.RWMutex
}

type meh struct {
	io.Reader
}

func (m meh) Close() error {
	return nil
}

//GetFrameJPEG returns an io.ReadCloser with a frame in JPEG format
func (l *Local) GetFrameJPEG() (io.ReadCloser, error) {
	l.lck.RLock()
	l.done.RLock()
	defer l.lck.RUnlock()
	defer l.done.RUnlock()
	if l.closed {
		return nil, errors.New("camera closed")
	}
	return meh{bytes.NewBuffer(l.curframe)}, nil
}

//Close a camera.Local
func (l *Local) Close() error {
	if l.closed {
		return nil
	}
	l.toclose = true
	l.done.Lock()
	l.done.Unlock()
	if !l.closed {
		return errors.New("Not closed????")
	}
	return nil
}

func fourcc(str string) webcam.PixelFormat {
	lc := []rune(str)
	if len(lc) != 4 {
		panic(fmt.Errorf("four letter code is not four letters (got %d)", len(lc)))
	}
	return webcam.PixelFormat(
		uint32(lc[0]) |
			(uint32(lc[1]) << 8) |
			(uint32(lc[2]) << 16) |
			(uint32(lc[3]) << 24),
	)
}

//NewLocalCam sets up a new local camera
func NewLocalCam(path string) (*Local, error) {
	l := new(Local)
	//open the device
	w, err := webcam.Open(path)
	if err != nil {
		return nil, err
	}
	//find pixel format
	fmts := w.GetSupportedFormats()
	pxfmt := webcam.PixelFormat(0)
	for _, f := range []webcam.PixelFormat{fourcc("MJPG"), fourcc("JPEG")} {
		if fmts[f] != "" {
			pxfmt = f
		}
	}
	if pxfmt == 0 {
		return nil, errors.New("No supported pixel format detected")
	}
	//find best resolution
	var resolution webcam.FrameSize
	var best uint64
	for _, sz := range w.GetSupportedFrameSizes(pxfmt) {
		pxn := uint64(sz.MaxHeight) * uint64(sz.MaxWidth)
		if pxn > best {
			resolution = sz
			best = pxn
		}
	}
	//configure camera
	_, _, _, err = w.SetImageFormat(pxfmt, resolution.MaxWidth, resolution.MaxHeight)
	if err != nil {
		if e := w.Close(); e != nil {
			panic(e) //if we cannot close, we have a resource leak
		}
		return nil, err
	}
	err = w.StartStreaming()
	if err != nil {
		if e := w.Close(); e != nil {
			panic(e) //if we cannot close, we have a resource leak
		}
		return nil, err
	}
	err = w.WaitForFrame(1)
	if err != nil {
		panic(err)
	}
	l.curframe, err = w.ReadFrame()
	if err != nil {
		panic(err)
	}
	l.done.Lock()
	go func() {
		defer l.done.Unlock()
		for !l.toclose {
			err := w.WaitForFrame(1)
			if err != nil {
				panic(err)
			}
			dat, err := w.ReadFrame()
			if err != nil {
				panic(err)
			}
			l.lck.Lock()
			l.curframe = dat
			l.lck.Unlock()
		}
		l.closed = true
		err = w.Close()
		if err != nil {
			panic(err)
		}
	}()
	return l, nil
}
