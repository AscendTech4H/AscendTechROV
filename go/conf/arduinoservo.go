package conf

import (
	"errors"

	"github.com/AscendTech4H/bracketconf"
)

//ArduinoServo is an Arduino servo pinout
type ArduinoServo struct {
	Name       string
	ControlPin int
}

var arduinoServoDirectiveProcessor = bracketconf.NewDirectiveProcessor(
	nameDirective,
	bracketconf.Directive{Name: "controlpin", Callback: func(object interface{}, ans ...bracketconf.ASTNode) {
		if len(ans) == 1 {
			object.(*ArduinoServo).ControlPin = ans[0].Int()
		} else if len(ans) == 0 {
			panic(errors.New("Enable directive has no arguments"))
		} else {
			panic(bracketconf.ConfErr{Pos: ans[0].Position(), Err: errors.New("Enable directive has too many arguments")})
		}
	}},
)

var arduinoServoDirective = bracketconf.Directive{Name: "ardservo", Callback: func(object interface{}, ans ...bracketconf.ASTNode) {
	ard := object.(*Arduino)
	switch len(ans) {
	case 0:
		panic(errors.New("Servo directive has no arguments"))
	case 1:
		if !ans[0].IsBracket() {
			panic(bracketconf.ConfErr{Pos: ans[0].Position(), Err: errors.New("Single servo argument must be a bracket")})
		}
		as := ArduinoServo{}
		ans[0].Evaluate(&as, arduinoServoDirectiveProcessor)
		ard.Servos = append(ard.Servos, &as)
		return
	case 2:
		ard.Servos = append(ard.Servos, &ArduinoServo{
			Name:       ans[0].Text(),
			ControlPin: ans[1].Int(),
		})
		return
	}
	for _, v := range ans {
		if v.IsBracket() {
			panic(bracketconf.ConfErr{Pos: ans[0].Position(), Err: errors.New("Invalid syntax: cannot mix bracket and non-bracket syntax for an arduino servo directive")})
		}
	}
	panic(bracketconf.ConfErr{Pos: ans[0].Position(), Err: errors.New("Invalid non-bracket syntax: must have 2 arguments for non-bracket syntax")})
}}
