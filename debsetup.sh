installstuff () {
	echo "deb http://llvm.org/apt/jessie/     llvm-toolchain-jessie-3.8 main" >  /etc/apt/sources.list.d/llvm.list
	apt-get install arduino arduino-mk llvm-3.8
	curl https://github.com/dibyendumajumdar/ravi/releases/download/0.15.1/ravi-ubuntu15-x64-llvm38.tar.gz | tar -xvf - ravi-ubuntu15-x64-llvm38/bin/ ravi-ubuntu15-x64-llvm38/lib/ -C /
}
echo "We need to install a few things for this to work"
if [ "$(id -u)" != "0" ]; then
	export -f installstuff
	su -c installstuff
else
	installstuff
fi
