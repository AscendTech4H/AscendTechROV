package can

import (
	"../commander"
)

type canSender struct {
	bus CAN
}

func (s canSender) Send(c commander.Command) {
	m := Message{}
	m.Data = commander.Serialize(c)
	m.Sender = s.bus.SenderID
	s.bus.SendMessage(m)
}

//AsSender returns a command sender for this CAN bus
func (c CAN) AsSender() commander.Sender {
	s := canSender{}
	s.bus = c
	return s
}
