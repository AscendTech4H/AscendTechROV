//Package motor implements a motor interface
package motor

//Type is the motor type
type Type uint8

//Types of motor
const (
	DC Type = iota
	Servo
	Stepper
)

//Motor is an interface to be implemented for motor drivers
type Motor interface {
	GetMotorType() Type
	Set(uint8)
	State() uint8
}

//TypeName gets the textual representation of the motor type
func (m Type) TypeName() string {
	switch m {
	case DC:
		return "DC"
	case Servo:
		return "servo"
	case Stepper:
		return "stepper"
	}
	panic("Uh oh")
}
