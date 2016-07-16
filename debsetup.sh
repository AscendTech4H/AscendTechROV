installstuff () {
	if [ ! -a /etc/apt/sources.list.d/llvm.list ] then
		echo "deb http://llvm.org/apt/jessie/     llvm-toolchain-jessie-3.8 main" >  /etc/apt/sources.list.d/llvm.list
	fi
	apt-get update
	apt-get install arduino llvm-3.8
	RAVILOC=$(mktemp -d)
	echo "Making ravi .deb package. . . "
	curl https://github.com/dibyendumajumdar/ravi/releases/download/0.15.1/ravi-ubuntu15-x64-llvm38.tar.gz | tar -xvf - ravi-ubuntu15-x64-llvm38/bin/ ravi-ubuntu15-x64-llvm38/lib/ -C $RAVILOC
	mv $RAVILOC/ravi-ubuntu15-x64-llvm38/bin $RAVILOC/bin
	mv $RAVILOC/ravi-ubuntu15-x64-llvm38/lib $RAVILOC/lib
	rm -rf $RAVILOC/ravi-ubuntu15-x64-llvm38
	mkdir $RAVILOC/debian
	(cd $RAVILOC; dch --create -v 0.15.1-1 --package ravi)
	echo 9 > $RAVILOC/debian/compat
	echo "Package: ravi" > $RAVILOC/debian/control
	echo "Version: 0.15.1" >> $RAVILOC/debian/control
	echo "Architecture: amd64" >> $RAVILOC/debian/control
	echo "Depends: llvm-3.8" >> $RAVILOC/debian/control
	echo "Description: Ravi is a derivative/dialect of Lua 5.3 with limited optional static typing and an LLVM powered JIT compiler." >> $RAVILOC/debian/control
	curl https://raw.githubusercontent.com/dibyendumajumdar/ravi/master/LICENSE > $RAVILOC/debian/copyright
	echo "#!/usr/bin/make -f" > $RAVILOC/debian/rules
	echo "%:" >> $RAVILOC/debian/rules
	echo "	dh $@" >> $RAVILOC/debian/rules
	echo "" >> $RAVILOC/debian/rules
	echo "override_dh_auto_install:" >> $RAVILOC/debian/rules
	echo "$(MAKE) DESTDIR=$$(pwd)/debian/hithere prefix=/usr install" >> $RAVILOC/debian/rules
	echo "bin" > $RAVILOC/debian/ravi.dirs
	echo "lib" >> $RAVILOC/debian/ravi.dirs
	mkdir $RAVILOC/ravi-0.15.1
	mv $RAVILOC/bin $RAVILOC/ravi-0.15.1/bin
	mv $RAVILOC/lib $RAVILOC/ravi-0.15.1/lib
	mv $RAVILOC/debian $RAVILOC/ravi-0.15.1/debian
	(cd $RAVILOC; debuild -us -uc)
	echo "Ravi package built. Installing. . . "
	dpkg -i $RAVILOC/ravi_0.15.1-1_amd64.deb
	rm -rf $RAVILOC
	echo "Done!"
}
echo "We need to install a few things for this to work"
if [ "$(id -u)" != "0" ]; then
	export -f installstuff
	su -c installstuff
else
	installstuff
fi
