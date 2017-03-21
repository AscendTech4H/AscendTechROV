package commander

type commandBase struct {
	id uint8
}

func (c commandBase) ID() uint8 {
	return c.id
}

type motorCMD struct {
	commandBase
	motor, pwm uint8
}

func (c motorCMD) Arguments() []byte {
	return []byte([]uint8{c.motor, c.pwm})
}

//SetMotor creates a motor set command
func SetMotor(motor uint8, pwm uint8) Command {
	c := motorCMD{}
	c.id = 0
	c.motor = motor
	c.pwm = pwm
	return c
}

type laserCMD struct {
	commandBase
	state bool
}

func (c laserCMD) Arguments() []byte {
	s := byte(0)
	if c.state {
		s = 1
	}
	return []byte{s}
}

//SetLaser generates a command to set the laser state
func SetLaser(state bool) Command {
	c := laserCMD{}
	c.id = 1
	c.state = state
	return c
}
