//System to work with the can bus
package can

import (
	"../util"
	"github.com/brutella/can"
	"math/rand"
)

//CAN bus
type CAN struct {
	MessageChan chan Message //Read messages from this
	bus         *can.Bus
	SenderID    uint8
}

//SetupCAN sets up a CAN bus
func SetupCAN(id uint8) (c CAN) {
	c.MessageChan = make(chan Message)
	c.SenderID = id
	bus, err := can.NewBusForInterfaceWithName("can0")
	util.UhOh(err)
	bus.SubscribeFunc(func(f can.Frame) {
		var m Message
		m.Sender = f.Data[0]
		l := f.Data[1]
		m.Data = f.Data[2 : l+2]
		c.MessageChan <- m
	})
	bus.ConnectAndPublish()
	return
}

//Send a message
func (c CAN) SendMessage(m Message) {
	var f can.Frame
	f.ID = uint32(rand.Int31n((1 << 11) - 1)) //Create random ID
	f.Length = uint8(len(m.Data))             //Set length
	for i, v := range m.Data {                //Copy data to frame
		f.Data[i] = v
	}
	c.bus.Publish(f) //Send frame
}

//Message object
type Message struct {
	Sender uint8
	Data   []byte
}
