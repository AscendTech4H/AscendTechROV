//Package can - a system to work with the can bus
package can

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"log"
	"os/exec"
	"sync"

	"../commander"
	"../debug"
	"../startup"
	"../util"

	"github.com/huin/goserial"
)

//CAN bus
type CAN struct {
	bus  io.ReadWriteCloser
	lck  sync.Mutex
	scan *bufio.Scanner
}

//SetupCAN sets up a CAN bus
func SetupCAN(port string) *CAN {
	c := new(CAN)
	bus, err := goserial.OpenPort(&goserial.Config{
		Name: port,
		Baud: 115200,
	})
	util.UhOh(err)
	debug.VLog("Start cat")
	n := exec.Command("/bin/cat", port)
	c.bus = bus
	o, err := n.StdoutPipe()
	util.UhOh(err)
	util.UhOh(n.Start())
	c.scan = bufio.NewScanner(o)
	debug.VLog("Buffing")
	c.scan.Scan()
	log.Println(c.scan.Text())
	return c
}

//SendMessage sends a message
func (c *CAN) SendMessage(m Message) {
	c.lck.Lock()
	defer c.lck.Unlock()
	for v := range m {
		debug.VLog(fmt.Sprintf("%d", v))
	}
	_, err := c.bus.Write([]byte(m))
	if !c.scan.Scan() {
		panic("It no wirk")
	}
	log.Println(c.scan.Text())
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
		flag.StringVar(&canName, "can", "/dev/ttyACM0", "Can bus arduino port (default: /dev/ttyUSB0)")
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
