//Package can - a system to work with the can bus
package can

import (
	"bufio"
	"flag"
	"sync"

	"../commander"
	"../startup"
	"../util"

	"github.com/tarm/serial"
)

//CAN bus
type CAN struct {
	bus  *serial.Port
	lck  sync.Mutex
	scan *bufio.Scanner
}

//SetupCAN sets up a CAN bus
func SetupCAN(port string) *CAN {
	c := new(CAN)
	bus, err := serial.OpenPort(&serial.Config{
		Name: port,
		Baud: 115200,
	})
	util.UhOh(err)
	c.bus = bus
	c.scan = bufio.NewScanner(bus)
	return c
}

//SendMessage sends a message
func (c *CAN) SendMessage(m Message) {
	c.lck.Lock()
	defer c.lck.Unlock()
	_, err := c.bus.Write([]byte(m))
	util.UhOh(err)
}

//Message object
type Message []byte

//Args
var canName string

//Bus is the main CAN bus
var Bus *CAN

//Sender is the CAN command sender
var Sender commander.Sender

//NoCAN says if CAN is disabled
var NoCAN bool

func init() {
	startup.NewTask(1, func() error { //Set up can flag parsing
		flag.StringVar(&canName, "can", "can0", "Can bus (default: can0)")
		flag.BoolVar(&NoCAN, "nocan", false, "Whether can is disabled")
		return nil
	})
	startup.NewTask(100, func() error {
		if !NoCAN {
			Bus = SetupCAN(canName)
			Sender = Bus.AsSender()
		}
		return nil
	})
}
