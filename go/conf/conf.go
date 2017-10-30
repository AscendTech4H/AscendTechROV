package conf

import (
	"net/url"

	"../accelerometer"
	"github.com/AscendTech4H/AscendTechROV/go/motor"
	"github.com/AscendTech4H/bracketconf"
)

func init() {
	bracketconf.NewDirectiveProcessor()
}

//Robot is the thing
type Robot struct {
	Motors  map[string]motor.Motor
	Accel   map[string]accelerometer.Accel
	Cameras map[string]interface{}
	Startup []func(*Robot) error
}

//Arduino is an Arduino configuration
type Arduino struct {
	Serial       string
	Motors       []*ArduinoMotor
	Servos       []*ArduinoServo
	ArduinoAccel []*ArduinoAccel
}

//ArduinoMotor is a set of pins for an Arduino motor
type ArduinoMotor struct {
	Name      string
	Enable    int
	Direction [2]int
	PWM       int
}

//ArduinoServo is an Arduino servo pinout
type ArduinoServo struct {
	Name       string
	ControlPin int
}

//ArduinoAccel is an accelerometer configuration
type ArduinoAccel struct {
	Name string
	Addr []byte
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

type Control struct {
}

type Vertical struct {
	Motors []*MotAng
}

//MotAng is a motor and an angle
type MotAng struct {
	Motor string
	Ang   float64
}
