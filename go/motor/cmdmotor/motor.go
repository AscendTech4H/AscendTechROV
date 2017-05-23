//Package cmdmotor creates motor interfaces with a commander.Sender
package cmdmotor

import (
	"fmt"

	"github.com/AscendTech4H/AscendTechROV/go/motor"
	"github.com/AscendTech4H/AscendTechROV/go/commander"
	"github.com/AscendTech4H/AscendTechROV/go/debug"
)

type cmdmotor struct {
	sender    commander.Sender
	index     uint8
	motorType motor.Type
	state     uint8
}

func (m *cmdmotor) GetMotorType() motor.Type {
	return m.motorType
}

func (m *cmdmotor) Set(speed uint8) {
	debug.VLog(fmt.Sprintf("Set motor %d to %d", m.index, speed))
	if m.state != speed {
		m.sender.Send(commander.SetMotor(m.index, speed))
	}
	m.state = speed
}

func (m *cmdmotor) State() uint8 {
	return m.state
}

//Motor creates a new motor with the command sender
func Motor(s commander.Sender, index uint8, t motor.Type) motor.Motor {
	m := new(cmdmotor)
	m.sender = s
	m.index = index
	m.motorType = t
	m.state = 0
	return m
}
