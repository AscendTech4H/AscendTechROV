package conf

import (
	"github.com/AscendTech4H/AscendTechROV/go/accelerometer"
	"github.com/AscendTech4H/AscendTechROV/go/motor"
)

//Robot is the thing
type Robot struct {
	Motors  map[string]motor.Motor
	Accel   map[string]accelerometer.Accel
	Cameras map[string]interface{}
	Startup []func(*Robot) error
}

//TODO: this
