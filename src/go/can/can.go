//System to work with the can bus using sys/unix
package can

import (
	"golang.org/x/sys/unix"
	"os"
	"github.com/howeyc/crc16"
	"encoding/binary"
)

var fd int
func init() {
	f,err := unix.Socket(AF_CAN, SOCK_RAW, CAN_RAW)		//Make a SocketCAN socket
	panic(err)					//Crash if error
	err = unix.Bind(f,&SockaddrCAN{Ifindex: 0})		//Bind it to all CAN interfaces
	panic(err)					//Crash if error
}

func crc(dat []byte) (out []byte) {
	out = make(byte,2)
	binary.LittleEdian.PutUint16(out,crc16.Checksum(dat,crc16.makeTable(0x4599)))
}
