//Package communication - a system to work with the serial bus
package communication

import (
	"bufio"
	"flag"
	"io"
	"io/ioutil"
	"log"
	"sync"
	"time"

	"github.com/AscendTech4H/AscendTechROV/go/commander"
	"github.com/AscendTech4H/AscendTechROV/go/startup"
	"github.com/AscendTech4H/AscendTechROV/go/util"

	"github.com/huin/goserial"
)

type shmeh struct {
	s chan []byte
	b []byte
}

func (s *shmeh) Read(b []byte) (n int, err error) {
	var d []byte
	if s.b == nil {
		d = <-s.s
	} else {
		d = s.b
		s.b = nil
	}
	for i, v := range d {
		if i > len(b) {
			s.b = d[len(b):]
			return len(b), nil
		}
		b[i] = v
	}
	return len(d), nil
}

//CAN bus
type CAN struct {
	bus  io.ReadWriteCloser
	lck  sync.Mutex
	rch  *shmeh
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
	c.bus = bus
	c.rch = new(shmeh)
	c.rch.s = make(chan []byte)
	go func() {
		t := time.Tick(time.Second / 10)
		for {
			<-t
			b, e := ioutil.ReadFile(port)
			util.UhOh(e)
			if len(b) > 0 {
				c.rch.s <- b
			}
		}
	}()
	c.scan = bufio.NewScanner(c.rch)
	c.scan.Scan()
	log.Println(c.scan.Text())
	return c
}

//SendMessage sends a message
func (c *CAN) SendMessage(m Message) {
	c.lck.Lock()
	defer c.lck.Unlock()
	log.Println(m)
	/*s := ""
	for _, v := range m {
		s += fmt.Sprintln(uint(v) + 1)
	}
	log.Println(s)*/
	_, err := c.bus.Write([]byte(m))
	util.UhOh(err)
	for i := 0; i < len(m); i++ {
		c.scan.Scan()
		log.Println(c.scan.Text())
	}
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
		flag.StringVar(&canName, "can", "/dev/ttyACM0", "Can bus arduino port (default: /dev/ttyACM0)")
		flag.BoolVar(&NoCAN, "nocan", false, "Whether can is disabled")
		return nil
	})
	startup.NewTask(20, func() error {
		if !NoCAN {
			Bus = SetupCAN(canName)
			Sender = Bus.AsSender()
		}
		return nil
	})
}
