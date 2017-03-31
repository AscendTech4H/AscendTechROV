all: buildmain

buildmain: installdeps
	$(MAKE) -C main all

clean:
	$(MAKE) -C main clean

installdeps:
	go get github.com/blackjack/webcam github.com/brutella/can github.com/gorilla/websocket

test:
	$(MAKE) buildmain
