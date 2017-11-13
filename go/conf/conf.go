package conf

import (
	"net/url"
	"strconv"

	"../accelerometer"
	"github.com/AscendTech4H/AscendTechROV/go/motor"
	"github.com/AscendTech4H/bracketconf"
)

var dirProcessor bracketconf.DirectiveProcessor

type valueList struct {
	arr []interface{}
}

func (values *valueList) add(val string) {
	//try converting to int or float, then append whatever you get
	i, err := strconv.Atoi(val)
	if err != nil {
		f, err := strconv.ParseFloat(val, 64)
		if err != nil {
			values.arr = append(values.arr, val)
		} else {
			values.arr = append(values.arr, f)
		}
	} else {
		values.arr = append(values.arr, i)
	}
}
func (values *valueList) parseTree(n bracketconf.ASTNode) {
	switch {
	case n.IsDir():
		fallthrough
	case n.IsBracket():
		n.Evaluate(values, dirProcessor)
	case n.IsArr():
		n.ForEach(func(_ int, v bracketconf.ASTNode) {
			values.parseTree(v)
		})
	default:
		values.add(n.Text())
	}
}

func init() {
	arduinoMotorDirective := bracketconf.Directive{Name: "ardmotor", Callback: func(object interface{}, ans ...bracketconf.ASTNode) {
		values := valueList{}
		for _, n := range ans {
			values.parseTree(n)
		}
	}}
	dirProcessor = bracketconf.NewDirectiveProcessor(arduinoMotorDirective)
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
