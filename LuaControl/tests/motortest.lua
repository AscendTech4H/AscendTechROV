require "core/commands.lua"

arguments={...}
serport=arguments[1]
if serport==nil then
	print("Please enter a serial port to run the motor test on.")
	print("Presently connected serial ports: ")
	for i,v in ipairs(getSerialPorts()) do
		print("\t",v)
	end
	return 1
end
initPort(serport)
i=0
while i<8 do
	updateMotor(i,0)
end
i=0
updateTlc()
while true do
	if i==8 then
		i=0
	end
	print("Testing motor "..i)
	updateMotor(i,-127)
	updateTlc()
	os.execute("sleep 1")
	updateMotor(i,127)
	updateTlc()
	os.execute("sleep 1")
	updateMotor(i,0)
	i=i+1
end
