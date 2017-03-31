all: buildout

buildstatic:
	$(MAKE) -C static all

buildgo:
	$(MAKE) -C go all

buildmotordriver:
	$(MAKE) -C MotorDriver all

clean: cleanstatic cleango cleanout cleanmotordriver

cleanstatic:
	$(MAKE) -C static clean

cleango:
	$(MAKE) -C go clean

cleanmotordriver:
	$(MAKE) -C MotorDriver clean

cleanout:
	$(MAKE) -C out clean

test:
	$(MAKE) -C static test
	$(MAKE) -C go test

buildout: buildstatic buildgo buildmotordriver
	$(MAKE) -C out all
