local arguments=arg
if #arguments<2 then
	for l in io.lines(debug.getinfo(1).source:match("@?(.*/)").."help.txt") do
		print(l)
	end
	return
end
print("Making Folder")
os.execute("mkdir -p "..arguments[1])
print("Opening Makefile. . . ")
f=io.open(arguments[1].."/Makefile","w+")
assert(f)
f:write("export BOARD_TAG="..arguments[2].."\n")
f:write("export ARDUINO_LIBS=\"")
local i: interger
i=2
while i < #arguments do
	i=i+1
	if i>3 then f:write(" ") end
	f:write(arguments[i])
end
f:write("\"\n")
for l in io.lines(debug.getinfo(1).source:match("@?(.*/)").."template.mk") do
	f:write(l.."\n")
end
f:close()
print("Done!")
