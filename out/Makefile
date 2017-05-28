GOPATH = $(realpath .)/gopath

all: topSide.o camrelay.o static

topSide.o: prepgo
	go build -v -o topSide.o github.com/AscendTech4H/AscendTechROV/go/main

camrelay.o: prepgo
	go build -v -o camrelay.o github.com/AscendTech4H/AscendTechROV/go/camrelay

static:
	cp -r ../static static
	rm static/Makefile static/.gitignore

clean:
	rm -rf topSide.o static camrelay.o gopath

prepgo: loadgo godeps

loadgo: gopath
	cp -r ../go gopath/src/github.com/AscendTech4H/AscendTechROV/go

gopath:
	mkdir gopath
	mkdir -p gopath/src/github.com/AscendTech4H/AscendTechROV

godeps: gopath
	go get github.com/blackjack/webcam github.com/gorilla/websocket github.com/huin/goserial

makeall:
	$(MAKE) -C ..

cleanall:
	$(MAKE) -C .. clean

rebuildall:
	$(MAKE) cleanall
	$(MAKE) makeall
