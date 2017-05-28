package communication

import (
	"github.com/AscendTech4H/AscendTechROV/go/commander"
)

type serialSender struct {
	bus *Serial
}

func (s serialSender) Send(c commander.Command) {
	s.bus.SendMessage(Message(commander.Serialize(c)))
}

//AsSender returns a command sender for this CAN bus
func (c *Serial) AsSender() commander.Sender {
	s := serialSender{}
	s.bus = c
	return s
}
