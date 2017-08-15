GOPATH = $(shell pwd)/gopath

all: topSide camRelay static

topSide: prepgo
	go build -v -o topSide github.com/AscendTech4H/AscendTechROV/go/main

camRelay: prepgo
	go build -v -o camRelay github.com/AscendTech4H/AscendTechROV/go/camrelay

gotest: prepgo
	go test github.com/AscendTech4H/AscendTechROV/go/controlDriver

static:
	cp -r ../static static
	rm static/Makefile static/.gitignore

clean:
	rm -rf topSide static camRelay gopath

prepgo: loadgo godeps

loadgo: gopath
	$(MAKE) $(GOPATH)/src/github.com/AscendTech4H/AscendTechROV/go

$(GOPATH)/src/github.com/AscendTech4H/AscendTechROV/go:
	ln -s $(shell pwd)/../go $(GOPATH)/src/github.com/AscendTech4H/AscendTechROV/go

gopath:
	mkdir gopath
	mkdir -p $(GOPATH)/src/github.com/AscendTech4H/AscendTechROV

godeps: gopath
	go get github.com/blackjack/webcam github.com/gorilla/websocket github.com/huin/goserial

makeall:
	$(MAKE) -C ..

cleanall:
	$(MAKE) -C .. clean

rebuildall:
	$(MAKE) cleanall
	$(MAKE) makeall

ci: all gotest
