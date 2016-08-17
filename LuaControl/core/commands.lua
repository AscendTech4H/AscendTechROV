require "serial.lua"

function updateMotor(motor,value)
	sendByte(0)
	sendByte(motor)
	sendByte(value)
end

function updateTlc()
	sendByte(1)
end

function mpuRead()
	sendByte(2)
end

function sendData()
	sendByte(3)
end

function updateServo(servo,value)
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
