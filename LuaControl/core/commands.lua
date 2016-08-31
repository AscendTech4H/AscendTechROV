require "serial.lua"


function updateMotor(motor: interger,value: number)
	sendByte(0)
	sendByte(motor)
	sendByte(127-(value*127))
end

function updateSensor(sensor: interger)
	sendByte(1)
	sendByte(sensor)
end

function updateTlc()
	updateSensor(0)
end

function mpuRead()
	updateSensor(1)
end

function sendData()
	sendByte(2)
end

function updateServo(servo: interger,value: interger)
	sendByte(4)
	sendByte(servo)
	sendByte(value)
end

function bmpRead()
	sendByte(5)
end

function readExTemp()
	sendByte(6)
end
