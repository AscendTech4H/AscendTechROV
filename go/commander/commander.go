//Package commander implements and sends commands
package commander

import "log"

//Command is an interface representing a sendable command
type Command interface {
	Arguments() []byte
	ID() uint8
}

//Serialize serializes a command
func Serialize(c Command) []byte {
	arg := c.Arguments()
	log.Println(arg)
	o := make([]byte, len(arg)+1)
	o[0] = c.ID()
	for i, v := range arg {
		o[i+1] = v
	}
	return o
}

//Sender repreesents anything that can send commands
type Sender interface {
	Send(Command)
}
