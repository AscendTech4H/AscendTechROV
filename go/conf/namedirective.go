package conf

import (
	"errors"

	"github.com/AscendTech4H/bracketconf"
)

var nameDirective = bracketconf.Directive{Name: "name", Callback: func(object interface{}, ans ...bracketconf.ASTNode) {
	if len(ans) == 1 {
		switch o := object.(type) {
		case *ArduinoMotor:
		case *ArduinoServo:
		case *ArduinoAccel: //to add more possible classes to name, add them at the end of this "case" chain
			o.Name = ans[0].Text()
		}
	} else if len(ans) == 0 {
		panic(errors.New("Name directive has no arguments"))
	} else {
		panic(bracketconf.ConfErr{Pos: ans[0].Position(), Err: errors.New("Name directive has too many arguments")})
	}
}}
