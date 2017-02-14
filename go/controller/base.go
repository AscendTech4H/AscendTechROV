// Controller interface for robot
package controller


type Controller interface {
	GetButtons() map[string]Button

}



type Button interface {

	State() uint8

}



type Joystick interface {

	State() (uint8,uint8)
}
