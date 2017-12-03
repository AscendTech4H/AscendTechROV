package conf

import (
	"errors"

	"github.com/AscendTech4H/bracketconf"
)

//ArduinoAccel is an accelerometer configuration
type ArduinoAccel struct {
	Name string
	Addr uint8
}

var arduinoAccelDirectiveProcessor = bracketconf.NewDirectiveProcessor(
	nameDirective,
	bracketconf.Directive{Name: "addr", Callback: func(object interface{}, ans ...bracketconf.ASTNode) {
		if len(ans) == 1 {
			i := ans[0].Int()
			if i < 0 || i > 127 {
				panic(bracketconf.ConfErr{Pos: ans[0].Position(), Err: errors.New("Address must be between 0 and 127 inclusive")})
			}
			object.(*ArduinoAccel).Addr = uint8(i)
		} else if len(ans) == 0 {
			panic(errors.New("Address directive has no arguments"))
		} else {
			panic(bracketconf.ConfErr{Pos: ans[0].Position(), Err: errors.New("Address directive has too many arguments")})
		}
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
		i := ans[1].Int()
		if i < 0 || i > 127 {
			panic(bracketconf.ConfErr{Pos: ans[1].Position(), Err: errors.New("Address must be between 0 and 127 inclusive")})
		}
		ard.ArduinoAccel = append(ard.ArduinoAccel, &ArduinoAccel{
			Name: ans[0].Text(),
			Addr: uint8(i)
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
