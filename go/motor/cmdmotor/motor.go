//Package cmdmotor creates motor interfaces with a commander.Sender
package cmdmotor

import (
	".."
	"../../commander"
)

type cmdmotor struct {
	sender    commander.Sender
	index     uint8
	motorType motor.Type
	state     uint8
}

func (m cmdmotor) GetMotorType() motor.Type {
	return m.motorType
}

func (m cmdmotor) Set(speed uint8) {
	m.sender.Send(commander.SetMotor(m.index, speed))
	m.state = speed
}

func (m cmdmotor) State() uint8 {
	return m.state
}

//Motor creates a new motor with the command sender
func Motor(s commander.Sender, index uint8, t motor.Type) motor.Motor {
	m := cmdmotor{}
	m.sender = s
	m.index = index
	m.motorType = t
	m.state = 0
	return m
}
