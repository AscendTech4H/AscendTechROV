package conf

import (
	"errors"

	"github.com/AscendTech4H/bracketconf"
)

var arduinoMotorDirectionDirectiveProcessor = bracketconf.NewDirectiveProcessor(
	bracketconf.Directive{Name: "x", Callback: func(object interface{}, ans ...bracketconf.ASTNode) {
		if len(ans) == 1 {
			object.(*ArduinoMotor).Direction[0] = ans[0].Int()
		} else {
			panic(errors.New("X directive needs 1 argument"))
		}
	}},
	bracketconf.Directive{Name: "y", Callback: func(object interface{}, ans ...bracketconf.ASTNode) {
		if len(ans) == 1 {
			object.(*ArduinoMotor).Direction[1] = ans[0].Int()
		} else {
			panic(errors.New("Y directive needs 1 argument"))
		}
	}},
)
var arduinoMotorDirectiveProcessor = bracketconf.NewDirectiveProcessor(
	bracketconf.Directive{Name: "name", Callback: func(object interface{}, ans ...bracketconf.ASTNode) {
		if len(ans) == 1 {
			object.(*ArduinoMotor).Name = ans[0].Text()
		} else {
			panic(errors.New("Name directive needs 1 argument"))
		}
	}},
	bracketconf.Directive{Name: "enable", Callback: func(object interface{}, ans ...bracketconf.ASTNode) {
		if len(ans) == 1 {
			object.(*ArduinoMotor).Enable = ans[0].Int()
		} else {
			panic(errors.New("Enable directive needs 1 argument"))
		}
	}},
	bracketconf.Directive{Name: "pwm", Callback: func(object interface{}, ans ...bracketconf.ASTNode) {
		if len(ans) == 1 {
			object.(*ArduinoMotor).PWM = ans[0].Int()
		} else {
			panic(errors.New("PWM directive needs 1 argument"))
		}
	}},
	bracketconf.Directive{Name: "direction", Callback: func(object interface{}, ans ...bracketconf.ASTNode) {
		if len(ans) == 1 {
			ans[0].Evaluate(object.(*ArduinoMotor), arduinoMotorDirectionDirectiveProcessor)
		} else if len(ans) == 2 {
			object.(*ArduinoMotor).Direction = [2]int{ans[0].Int(), ans[1].Int()}
		} else {
			panic(errors.New("Enable directive needs 1 bracket argument or 2 int arguments"))
		}
	}},
)
var arduinoMotorDirective = bracketconf.NewDirectiveProcessor(
	bracketconf.Directive{Name: "ardmotor", Callback: func(object interface{}, ans ...bracketconf.ASTNode) {
		ard := object.(*Arduino)
		switch len(ans) {
		case 0:
			panic(errors.New("Motor directive has no arguments"))
		case 1:
			if !ans[0].IsBracket() {
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
	}},
)
