package camera

import "testing"
import (
	"image/jpeg"
	"os"
	"time"
	"fmt"
)

func TestCamera(t *testing.T) {
	fmt.Println("Shmeh")
	SetupCam("/dev/video0")
	time.Sleep(time.Second*8)
	p := GetPic()
	f,err := os.Create("img.jpeg")
	if err!=nil {
		t.Error(err)
	}
	fmt.Println(f,p)
	err = jpeg.Encode(f,p,nil)
	if err!=nil {
		t.Error(err)
	}
}
