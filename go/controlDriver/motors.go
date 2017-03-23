//Package controlDriver issues commands based on controller input
package controlDriver

import (
	"../can"
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
//Note: put them in index order for iota to assign indexes (stating with 1)
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
	/*startup.NewTask(255, func() error {
		tick := time.NewTicker(5 * time.Second)
		go func() {
			for <-tick.C {
				r := controller.RobotState()
				//Do something when I am awake enough to know what I am doing
			}
		}()
		return nil
	})*/
}
