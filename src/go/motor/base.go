package motor

type MotorType uint8

const (
        DC MotorType = iota
        Servo
        Stepper
)

type Motor interface{
	GetMotorType() MotorType
	Set(uint8)
	State() uint8
}

func TypeName(m MotorType) string {
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
