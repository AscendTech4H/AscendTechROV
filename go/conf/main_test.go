package conf

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/AscendTech4H/bracketconf"
)

func TestArduino(t *testing.T) {
	tv, err := bracketconf.Parse(strings.NewReader(`
	arduino {
    ardmotor a 2 3 4 5;
    ardmotor {
      name b;
      enable 65;
      direction {
        x 0;
        y 1;
      };
      pwm 52;
    };
    ardmotor {
      name c;
      enable 65;
      pwm 52;
      direction 2 3;
    };
    ardservo a 4;
    ardservo {
      name b;
      controlpin 3;
    };
    ardaccel a 2;
    ardaccel {
      name b;
      addr 65;
    };
	};
	localcamera a asdf;
	localcamera {
		name b;
		devfile abc;
	};
	remotecamera a "http://google.com/";
	remotecamera {
		name b;
		address "http://yahoo.com/";
	};
	`), "testing.conf", bracketconf.NewDirectiveProcessor(arduinoDirective, localCameraDirective, remoteCameraDirective, motAngDirective), &FullList{[]interface{}{}})
	if err != nil {
		t.Fatal(err.Error())
	}
	dat, err := json.Marshal(tv.(*FullList))
	if err != nil {
		t.Fatal(err.Error())
	}
	if string(dat) != `{"Arr":[{"Serial":"","Motors":[{"Name":"a","Enable":2,"Direction":[3,4],"PWM":5},{"Name":"b","Enable":65,"Direction":[0,1],"PWM":52},{"Name":"c","Enable":65,"Direction":[2,3],"PWM":52}],"Servos":[{"Name":"a","ControlPin":4},{"Name":"b","ControlPin":3}],"ArduinoAccel":[{"Name":"a","Addr":2},{"Name":"b","Addr":65}]},{"Name":"a","DevFile":"asdf"},{"Name":"b","DevFile":"abc"},{"Name":"a","Address":{"Scheme":"http","Opaque":"","User":null,"Host":"google.com","Path":"/","RawPath":"","RawQuery":"","Fragment":""}},{"Name":"b","Address":{"Scheme":"http","Opaque":"","User":null,"Host":"yahoo.com","Path":"/","RawPath":"","RawQuery":"","Fragment":""}}]}` {
		t.Fatalf("Incorrect parse %s", string(dat))
	}
}
