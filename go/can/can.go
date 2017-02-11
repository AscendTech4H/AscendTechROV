//System to work with the can bus using sys/unix
package can

import (
	"encoding/binary"
	"github.com/howeyc/crc16"
	"golang.org/x/sys/unix"
	"os"
)

var fd int

func init() {
	f, err := unix.Socket(AF_CAN, SOCK_RAW, CAN_RAW) //Make a SocketCAN socket
	panic(err)                                       //Crash if error
	err = unix.Bind(f, &SockaddrCAN{Ifindex: 0})     //Bind it to all CAN interfaces
	panic(err)                                       //Crash if error
}