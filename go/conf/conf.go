package conf

import (
	"net/url"

	"../accelerometer"
	"github.com/AscendTech4H/AscendTechROV/go/motor"
	"github.com/AscendTech4H/bracketconf"
)

var dirProcessor bracketconf.DirectiveProcessor

func init() {
	dirProcessor = bracketconf.NewDirectiveProcessor(arduinoDirective)
}

//Robot is the thing
type Robot struct {
	Motors  map[string]motor.Motor
	Accel   map[string]accelerometer.Accel
	Cameras map[string]interface{}
	Startup []func(*Robot) error
}

//LocalCamera is a camera attached to the local computer
type LocalCamera struct {
	Name    string
	DevFile string
}

//RemoteCamera is a remotely connected camera
type RemoteCamera struct {
	Name    string
	Address *url.URL
}

//Control contains the control settings for each direction
type Control struct {
	Directions []*Direction
}

//Direction is the control settings for a certain direciton
type Direction struct {
	Axis          int
	Motors        []*MotAng
	Stabilization *ArduinoAccel
}

//MotAng is a motor and an angle
type MotAng struct {
	Motor string
	Ang   float64
}
