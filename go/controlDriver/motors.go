//Package controlDriver issues commands based on controller input
package controlDriver

import (
	"math"
	"time"

	"../can"
	"../controller"
	"../motor"
	"../motor/cmdmotor"
	"../startup"
)

//Add more motors when I know which they are
var robot struct {
	leftback, rightback   motor.Motor
	leftfront, rightfront motor.Motor
}

//Motor IDs
//Note: put them in index order for iota to assign indexes (stating with 0)
const (
	motorlb = iota
	motorrb
	motorlf
	motorrf
)

func init() {
	startup.NewTask(150, func() error {
		robot.leftback = cmdmotor.Motor(can.Sender, motorlb, motor.DC)
		robot.rightback = cmdmotor.Motor(can.Sender, motorrb, motor.DC)
		robot.leftfront = cmdmotor.Motor(can.Sender, motorlb, motor.DC)
		robot.rightfront = cmdmotor.Motor(can.Sender, motorrb, motor.DC)
		return nil
	})
	startup.NewTask(255, func() error {
		tick := time.NewTicker(5 * time.Second)
		go func() {
			for range tick.C {
				rob := controller.RobotState()
				l, r := motorCalcFwd(rob.Forward, rob.Turn)
				a := uint8(rangeMap(r, -127, 127, 0, 255))
				b := uint8(rangeMap(l, -127, 127, 0, 255))
				robot.rightfront.Set(a)
				robot.rightback.Set(a)
				robot.leftfront.Set(b)
				robot.leftback.Set(b)
			}
		}()
		return nil
	})
}

func rangeMap(in, inmin, inmax, outmin, outmax int) int {
	return (((in - inmin) * (outmax - outmin)) / (inmax - inmin)) + outmin
}

func motorCalcFwd(forward int, turn int) (l, r int) {
	ang := math.Atan(float64(forward) / float64(turn))
	mag := math.Sqrt(float64((forward * forward) + (turn * turn)))
	if turn < 0 {
		l = int(mag * math.Sin(ang))
		r = int(mag)
	} else if turn > 0 {
		l = int(mag)
		r = int(mag * math.Sin(ang))
	} else {
		l = int(mag)
		r = int(mag)
	}
	return
}
