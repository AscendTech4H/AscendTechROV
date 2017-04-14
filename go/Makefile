all: buildmain buildcamrelay

buildmain: installdeps
	$(MAKE) -C main all

buildcamrelay:
	$(MAKE) -C camrelay all

clean:
	$(MAKE) -C main clean
	$(MAKE) -C camrelay clean

installdeps:
	go get github.com/blackjack/webcam github.com/brutella/can github.com/gorilla/websocket github.com/huin/goserial

test:
	$(MAKE) buildmain buildcamrelay
