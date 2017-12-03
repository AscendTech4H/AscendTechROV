package conf

import (
	"errors"

	"github.com/AscendTech4H/bracketconf"
)

//Arduino is an Arduino configuration
type Arduino struct {
	Serial       string
	Motors       []*ArduinoMotor
	Servos       []*ArduinoServo
	ArduinoAccel []*ArduinoAccel
}

var arduinoSerialDirective = bracketconf.Directive{Name: "serial", Callback: func(object interface{}, ans ...bracketconf.ASTNode) {
	if len(ans) == 1 {
		object.(*Arduino).Serial = ans[0].Text()
	} else {
		panic(errors.New("Serial directive needs 1 argument"))
	}
}}

var arduinoDirective = bracketconf.Directive{Name: "arduino", Callback: func(object interface{}, ans ...bracketconf.ASTNode) {
	if len(ans) == 1 && ans[0].IsBracket() {
		ard := Arduino{}
		ans[0].Evaluate(&ard, bracketconf.NewDirectiveProcessor(
			arduinoSerialDirective,
			arduinoMotorDirective,
			arduinoServoDirective,
			arduinoAccelDirective,
		))
		object.(*FullList).add(ard)
	} else if len(ans) == 2 && !ans[0].IsBracket() && ans[1].IsBracket() {
		ard := Arduino{}
		ard.Serial = ans[0].Text()
		ans[1].Evaluate(&ard, bracketconf.NewDirectiveProcessor(
			arduinoMotorDirective,
			arduinoServoDirective,
			arduinoAccelDirective,
		))
		object.(*FullList).add(ard)
	} else {
		panic(errors.New("Arduino directive needs 1 bracket argument or a string argument and a bracket argument"))
	}
}}
