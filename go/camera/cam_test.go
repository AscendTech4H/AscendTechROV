package camera

import (
	"bytes"
	"errors"
	"image"
	"image/color"
	"image/jpeg"
	"io"
	"math"
	"math/rand"
	"net/http/httptest"
	"net/url"
	"sync"
	"testing"
)

type testCam struct {
	lck sync.RWMutex
	img image.Image
}

func (c *testCam) GetFrameJPEG() (io.ReadCloser, error) {
	c.lck.RLock()
	defer c.lck.RUnlock()
	buf := bytes.NewBuffer(nil)
	err := jpeg.Encode(buf, c.img, nil)
	if err != nil {
		return nil, err
	}
	return meh{buf}, nil
}

func (c *testCam) Close() error {
	c.lck.Lock()
	defer c.lck.Unlock()
	if c.img == nil {
		return errors.New("already closed")
	}
	c.img = nil
	return nil
}

func (c *testCam) check(r io.ReadCloser, t *testing.T) {
	defer r.Close()
	img, err := jpeg.Decode(r)
	if err != nil {
		t.Fatal(err)
	}
	diff := 0.0
	for x := 0; x < 1280; x++ {
		for y := 0; y < 720; y++ {
			r1, g1, b1, _ := img.At(x, y).RGBA()
			r2, g2, b2, _ := c.img.At(x, y).RGBA()
			diff += math.Abs(float64(r1) - float64(r2))
			diff += math.Abs(float64(g1) - float64(g2))
			diff += math.Abs(float64(b1) - float64(b2))
		}
	}
	diff /= float64(1280 * 720 * 3 * 0xffff)
	diff *= 100
	if diff > 25 {
		t.Fatalf("ERROR: large difference (%f%%)", diff)
	}
	t.Logf("Difference: %f%%", diff)
}

func genTestCam(seed int64) *testCam {
	r := rand.New(rand.NewSource(seed))
	img := image.NewNRGBA(image.Rect(0, 0, 1280, 720))
	for x := 0; x < 1280; x++ {
		for y := 0; y < 720; y++ {
			img.SetNRGBA(x, y, color.NRGBA{
				uint8(r.Intn(255)),
				uint8(r.Intn(255)),
				uint8(r.Intn(255)),
				255,
			})
		}
	}
	return &testCam{img: img}
}

func TestTestCam(t *testing.T) {
	tc := genTestCam(65)
	r, err := tc.GetFrameJPEG()
	if err != nil {
		t.Fatal(err)
	}
	tc.check(r, t)
}

func TestRemoteCam(t *testing.T) {
	tc := genTestCam(52)
	hts := httptest.NewServer(&CamHandler{tc})
	defer hts.Close()
	u, err := url.Parse(hts.URL)
	if err != nil {
		t.Fatal(err)
	}
	rcam := NewRemoteCam(u, hts.Client())
	r, err := rcam.GetFrameJPEG()
	if err != nil {
		t.Fatal(err)
	}
	tc.check(r, t)
	err = rcam.Close()
	if err != nil {
		t.Fatal(err)
	}
	err = rcam.Close()
	if err == nil {
		t.Fatal("camera.Remote did not throw double close error")
	}
}
