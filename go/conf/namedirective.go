package conf

import (
	"errors"

	"github.com/AscendTech4H/bracketconf"
)

var nameDirective = bracketconf.Directive{Name: "name", Callback: func(object interface{}, ans ...bracketconf.ASTNode) {
	if len(ans) == 1 {
		switch o := object.(type) {
		case *ArduinoMotor:
			o.Name = ans[0].Text()
		case *ArduinoServo:
			o.Name = ans[0].Text()
		case *ArduinoAccel:
			o.Name = ans[0].Text()
		case *LocalCamera:
			o.Name = ans[0].Text()
		case *RemoteCamera:
			o.Name = ans[0].Text()
		case *MotAng:
			o.Name = ans[0].Text()
		}
	} else if len(ans) == 0 {
		panic(errors.New("Name directive has no arguments"))
	} else {
		panic(bracketconf.ConfErr{Pos: ans[0].Position(), Err: errors.New("Name directive has too many arguments")})
	}
}}
