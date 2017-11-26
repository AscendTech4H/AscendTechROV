package conf

import (
	"errors"

	"github.com/AscendTech4H/bracketconf"
)

//ArduinoAccel is an accelerometer configuration
type ArduinoAccel struct {
	Name string
	Addr []byte
}

var arduinoAccelDirectiveProcessor = bracketconf.NewDirectiveProcessor(
	nameDirective,
	bracketconf.Directive{Name: "addr", Callback: func(object interface{}, ans ...bracketconf.ASTNode) {
		//TODO: address
	}},
)

var arduinoAccelDirective = bracketconf.Directive{Name: "ardaccel", Callback: func(object interface{}, ans ...bracketconf.ASTNode) {
	ard := object.(*Arduino)
	switch len(ans) {
	case 0:
		panic(errors.New("Acceleration directive has no arguments"))
	case 1:
		if !ans[0].IsBracket() {
		}
		aa := ArduinoAccel{}
		ans[0].Evaluate(&aa, arduinoMotorDirectiveProcessor)
		ard.ArduinoAccel = append(ard.ArduinoAccel, &aa)
		return
	case 2:
		ard.ArduinoAccel = append(ard.ArduinoAccel, &ArduinoAccel{
			Name: ans[0].Text(),
			//TODO: address
		})
		return
	}
	for _, v := range ans {
		if v.IsBracket() {
			panic(bracketconf.ConfErr{Pos: ans[0].Position(), Err: errors.New("Invalid syntax: cannot mix bracket and non-bracket syntax for an arduino acceleration directive")})
		}
	}
	panic(bracketconf.ConfErr{Pos: ans[0].Position(), Err: errors.New("Invalid non-bracket syntax: must have 2 arguments for non-bracket syntax")})
}}
