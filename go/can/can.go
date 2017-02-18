//System to work with the can bus
package can

import (
	"github.com/brutella/can"
)

func init() {
	bus, _ := can.NewBusForInterfaceWithName("can0")
	bus.ConnectAndPublish()
}
