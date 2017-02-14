//Motor interface
package motor

// MotorType is the motor type
type MotorType uint8

const (
        DC MotorType = iota
        Servo
        Stepper
)

// Motor is an interface to be implemented for motor drivers
type Motor interface{
	GetMotorType() MotorType
	Set(uint8)
	State() uint8
}


// TypeName gets the textual representation of the motor type
func (m MotorType) TypeName() string {
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
