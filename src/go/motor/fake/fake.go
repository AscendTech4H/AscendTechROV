package fake
import (
	"../../motor"
	"fmt"
)

type fakemotor struct{
	t motor.MotorType
	name string
	state uint8
}

func (f *fakemotor) Set(s uint8) {
	f.state=s
	fmt.Printf("Set motor %s to %u.\n",f.name,s)
}

func (f *fakemotor) State() uint8 {
	return f.state
}

func (f *fakemotor) GetMotorType() motor.MotorType {
	return f.t
}

func NewFake(name string, t motor.MotorType, istate uint8) motor.Motor {
	m := new(fakemotor)
	m.name=name
	m.t=t
	m.state=istate
	return m
}
