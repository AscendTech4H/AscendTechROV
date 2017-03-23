//Package can - a system to work with the can bus
package can

import (
	"flag"
	"math/rand"

	"../commander"
	"../startup"
	"../util"

	"github.com/brutella/can"
)

//CAN bus
type CAN struct {
	MessageChan chan Message //Read messages from this
	bus         *can.Bus
	SenderID    uint8
}

//SetupCAN sets up a CAN bus
func SetupCAN(id uint8, port string) (c CAN) {
	c.MessageChan = make(chan Message)
	c.SenderID = id
	bus, err := can.NewBusForInterfaceWithName(port)
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

//SendMessage sends a message
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

var canName string

//Bus is the main CAN bus
var Bus CAN

//Sender is the CAN command sender
var Sender commander.Sender

func init() {
	startup.NewTask(1, func() error { //Set up can flag parsing
		flag.StringVar(&canName, "can", "can0", "Can bus (default: can0)")
		return nil
	})
	startup.NewTask(200, func() error {
		Bus = SetupCAN(65, canName)
		Sender = Bus.AsSender()
		return nil
	})
}
