package accelerometer

//Vec3 is a 3D vector
type Vec3 struct {
	X, Y, Z int
}

//Reading is an accelerometer reading
type Reading struct {
	Gyro  Vec3
	Accel Vec3
}

//Accel is an accelerometer interface
type Accel interface {
	Read() (*Reading, error)
}
