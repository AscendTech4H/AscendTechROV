.PHONY: fetchluaofficial cleanluaversion makeluaofficial setuplua createdirs clean cleanbin all
fetchluaofficial:
		if [ ! -d lua ]; then curl https://www.lua.org/ftp/lua-5.3.3.tar.gz | gzip -d | tar -xvf -;mv lua-5.3.3 lua;fi

cleanluaversion:
	rm -rf lua

makeluaofficial: fetchluaofficial
	if [ ! -d bin ]; then $(MAKE) -C lua linux CFLAGS=' -O3'; $(MAKE) -C lua test; fi

createdirs:
	if [ ! -d bin ]; then mkdir bin; fi
	if [ ! -d bin/pip ]; then mkdir bin/pip; fi
setuplua: makeluaofficial createdirs
	cp lua/src/lua bin/lua
	cp lua/src/luac bin/luac

luashell:
	if [ ! -d bin ]; then make setuplua; fi
	./bin/lua

clean: cleanluaversion cleanbin

cleanbin:
	rm -rf bin

fetchsysdepends:
	@echo "We need to install a few things for this to work"
	su -c "aptitude install python-pip gcc-avr avrdude g++ libreadline-dev"

fetchplatformio:
	if [ ! -d bin/pip/platformio ]; then pip install platformio -t bin/pip; fi

setupplatformio: fetchplatformio
	bash platformioproject/setup

setup: fetchsysdepends setuplua setupplatformio
