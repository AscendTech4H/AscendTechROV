package communication

import (
	"github.com/AscendTech4H/AscendTechROV/go/commander"
)

type canSender struct {
	bus *CAN
}

func (s canSender) Send(c commander.Command) {
	s.bus.SendMessage(Message(commander.Serialize(c)))
}

//AsSender returns a command sender for this CAN bus
func (c *CAN) AsSender() commander.Sender {
	s := canSender{}
	s.bus = c
	return s
}
