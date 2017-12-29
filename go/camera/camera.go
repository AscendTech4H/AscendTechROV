//Package camera handles the camera processing
package camera

import (
	"io"
)

//Camera is a camera used for the robot
type Camera interface {
	GetFrameJPEG() (io.ReadCloser, error)
	Close() error
}
