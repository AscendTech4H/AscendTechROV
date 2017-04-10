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
	leftback, rightback         motor.Motor
	leftfront, rightfront       motor.Motor
	topleftfront, toprightfront motor.Motor
	topleftback, toprightback   motor.Motor
	claw                        struct {
		roll motor.Motor
		grab motor.Motor
	}
}

//Motor IDs
//Note: put them in index order for iota to assign indexes (stating with 0)
const (
	motorlb = iota
	motorrb
	motorlf
	motorrf
	motortlf
	motortrf
	motortlb
	motortrb
	motorroll
	motorgrab
)

func init() {
	startup.NewTask(150, func() error {
		robot.leftback = cmdmotor.Motor(can.Sender, motorlb, motor.DC)
		robot.rightback = cmdmotor.Motor(can.Sender, motorrb, motor.DC)
		robot.leftfront = cmdmotor.Motor(can.Sender, motorlb, motor.DC)
		robot.rightfront = cmdmotor.Motor(can.Sender, motorrb, motor.DC)

		robot.topleftfront = cmdmotor.Motor(can.Sender, motorlf, motor.DC)
		robot.toprightfront = cmdmotor.Motor(can.Sender, motorrf, motor.DC)
		robot.topleftback = cmdmotor.Motor(can.Sender, motorlb, motor.DC)
		robot.toprightback = cmdmotor.Motor(can.Sender, motorrb, motor.DC)

		robot.claw.roll = cmdmotor.Motor(can.Sender, motorroll, motor.DC)
		robot.claw.grab = cmdmotor.Motor(can.Sender, motorgrab, motor.Servo)
		return nil
	})
	startup.NewTask(255, func() error {
		tick := time.NewTicker(time.Second / 5)
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
				robot.topleftback.Set(uint8(rob.Up))
				robot.topleftfront.Set(uint8(rob.Up))
				robot.topleftback.Set(uint8(rob.Up))
				robot.toprightback.Set(uint8(rob.Up))
				c := uint8(0)
				switch rob.ClawTurn {
				case controller.CCW:
					c = 0
				case controller.CW:
					c = 255
				case controller.STOP:
					c = 127
				}
				robot.claw.roll.Set(c)
				if rob.Claw {
					robot.claw.grab.Set(180)
				} else {
					robot.claw.grab.Set(90)
				}
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
