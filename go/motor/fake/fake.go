//Package for implementing fake motors for testing purposes
package fake
import (
	"../../motor"
	"log"
)

type fakemotor struct{
	t motor.MotorType
	name string
	state uint8
}

func (f *fakemotor) Set(s uint8) {
	f.state=s
	log.Printf("Set %s motor %s to %d.\n",f.t.TypeName(),f.name,s)
}

func (f *fakemotor) State() uint8 {
	return f.state
}

func (f *fakemotor) GetMotorType() motor.MotorType {
	return f.t
}

// NewFake creates a new fake motor
func NewFake(name string, t motor.MotorType, istate uint8) motor.Motor {
	m := new(fakemotor)
	m.name=name
	m.t=t
	m.state=istate
	return m
}
