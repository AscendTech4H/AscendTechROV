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
