local f
function initPort(port)
	f=io.open("/dev/"..port,"rwb")
	if f==nil then
		error("Failed to open serial port "..port)
	end
end

function sendByte(b: interger)
	io.write(string.char(b))
end

function getSerialPorts(): table
	local i: interger = 1
	local ret: table = {}
	for l in io.popen("ls /dev | grep ttyACM"):lines() do
		ret[i]=l
		i=i+1
	end
	for l in io.popen("ls /dev | grep ttyUSB"):lines() do
		ret[i]=l
		i=i+1
	end
	ret.count=i-1
	return ret
end
