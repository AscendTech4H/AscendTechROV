GOPATH = $(shell pwd)/gopath
BINARIES = topside.o camrelay.o
GODEPS = github.com/blackjack/webcam github.com/gorilla/websocket github.com/huin/goserial

all: $(BINARIES) static

topside.o: prepgo
	go build -v -o topside.o github.com/AscendTech4H/AscendTechROV/go/main

camrelay.o: prepgo
	go build -v -o camrelay.o github.com/AscendTech4H/AscendTechROV/go/camrelay

gotest: prepgo
	go test github.com/AscendTech4H/AscendTechROV/go/controlDriver

static:
	cp -r ../static static
	rm static/Makefile static/.gitignore

clean:
	rm -rf static gopath $(BINARIES)

prepgo: loadgo godeps

loadgo: gopath
	$(MAKE) $(GOPATH)/src/github.com/AscendTech4H/AscendTechROV/go

$(GOPATH)/src/github.com/AscendTech4H/AscendTechROV/go:
	ln -s $(shell pwd)/../go $(GOPATH)/src/github.com/AscendTech4H/AscendTechROV/go

gopath:
	mkdir gopath
	mkdir -p $(GOPATH)/src/github.com/AscendTech4H/AscendTechROV

godeps: gopath
	go get $(GODEPS)

makeall:
	$(MAKE) -C ..

cleanall:
	$(MAKE) -C .. clean

rebuildall:
	$(MAKE) cleanall
	$(MAKE) makeall

ci: all gotest
