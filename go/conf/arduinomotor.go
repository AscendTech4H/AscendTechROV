package conf

import (
	"errors"

	"github.com/AscendTech4H/bracketconf"
)

//ArduinoMotor is a set of pins for an Arduino motor
type ArduinoMotor struct {
	Name      string
	Enable    int
	Direction [2]int
	PWM       int
}

var arduinoMotorDirectionDirectiveProcessor = bracketconf.NewDirectiveProcessor(
	bracketconf.Directive{Name: "x", Callback: func(object interface{}, ans ...bracketconf.ASTNode) {
		if len(ans) == 1 {
			object.(*ArduinoMotor).Direction[0] = ans[0].Int()
		} else if len(ans) == 0 {
			panic(errors.New("X directive has no arguments"))
		} else {
			panic(bracketconf.ConfErr{Pos: ans[0].Position(), Err: errors.New("X directive has too many arguments")})
		}
	}},
	bracketconf.Directive{Name: "y", Callback: func(object interface{}, ans ...bracketconf.ASTNode) {
		if len(ans) == 1 {
			object.(*ArduinoMotor).Direction[1] = ans[0].Int()
		} else if len(ans) == 0 {
			panic(errors.New("Y directive has no arguments"))
		} else {
			panic(bracketconf.ConfErr{Pos: ans[0].Position(), Err: errors.New("Y directive has too many arguments")})
		}
	}},
)
var arduinoMotorDirectiveProcessor = bracketconf.NewDirectiveProcessor(
	bracketconf.Directive{Name: "name", Callback: func(object interface{}, ans ...bracketconf.ASTNode) {
		if len(ans) == 1 {
			object.(*ArduinoMotor).Name = ans[0].Text()
		} else if len(ans) == 0 {
			panic(errors.New("Name directive has no arguments"))
		} else {
			panic(bracketconf.ConfErr{Pos: ans[0].Position(), Err: errors.New("Name directive has too many arguments")})
		}
	}},
	bracketconf.Directive{Name: "enable", Callback: func(object interface{}, ans ...bracketconf.ASTNode) {
		if len(ans) == 1 {
			object.(*ArduinoMotor).Enable = ans[0].Int()
		} else if len(ans) == 0 {
			panic(errors.New("Enable directive has no arguments"))
		} else {
			panic(bracketconf.ConfErr{Pos: ans[0].Position(), Err: errors.New("Enable directive has too many arguments")})
		}
	}},
	bracketconf.Directive{Name: "pwm", Callback: func(object interface{}, ans ...bracketconf.ASTNode) {
		if len(ans) == 1 {
			object.(*ArduinoMotor).PWM = ans[0].Int()
		} else if len(ans) == 0 {
			panic(errors.New("PWM directive has no arguments"))
		} else {
			panic(bracketconf.ConfErr{Pos: ans[0].Position(), Err: errors.New("PWM directive has too many arguments")})
		}
	}},
	bracketconf.Directive{Name: "direction", Callback: func(object interface{}, ans ...bracketconf.ASTNode) {
		if len(ans) == 1 && ans[0].IsBracket() {
			ans[0].Evaluate(object.(*ArduinoMotor), arduinoMotorDirectionDirectiveProcessor)
		} else if len(ans) == 2 {
			object.(*ArduinoMotor).Direction = [2]int{ans[0].Int(), ans[1].Int()}
		} else if len(ans) == 0 {
			panic(errors.New("Enable directive has no arguments"))
		} else {
			panic(bracketconf.ConfErr{Pos: ans[0].Position(), Err: errors.New("Enable directive needs 1 bracket argument or 2 int arguments")})
		}
	}},
)
var arduinoMotorDirective = bracketconf.Directive{Name: "ardmotor", Callback: func(object interface{}, ans ...bracketconf.ASTNode) {
	ard := object.(*Arduino)
	switch len(ans) {
	case 0:
		panic(errors.New("Motor directive has no arguments"))
	case 1:
		if !ans[0].IsBracket() {
			panic(bracketconf.ConfErr{Pos: ans[0].Position(), Err: errors.New("Motor directive's argument should be bracketed")})
		}
		am := ArduinoMotor{}
		ans[0].Evaluate(&am, arduinoMotorDirectiveProcessor)
		ard.Motors = append(ard.Motors, &am)
		return
	case 5:
		ard.Motors = append(ard.Motors, &ArduinoMotor{
			Name:      ans[0].Text(),
			Enable:    ans[1].Int(),
			Direction: [2]int{ans[2].Int(), ans[3].Int()},
			PWM:       ans[4].Int(),
		})
		return
	}
	for _, v := range ans {
		if v.IsBracket() {
			panic(bracketconf.ConfErr{Pos: ans[0].Position(), Err: errors.New("Invalid syntax: cannot mix bracket and non-bracket syntax for an arduino motor directive")})
		}
	}
	panic(bracketconf.ConfErr{Pos: ans[0].Position(), Err: errors.New("Invalid non-bracket syntax: must have 5 arguments for non-bracket syntax")})
}}
