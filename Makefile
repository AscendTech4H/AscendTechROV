all: buildout

buildstatic:
	$(MAKE) -C static all

buildgo:
	$(MAKE) -C go all

clean: cleanstatic cleanout

cleanstatic:
	$(MAKE) -C static clean

cleango:
	$(MAKE) -C go clean

cleanout:
	$(MAKE) -C out clean

test:
	$(MAKE) -C static test
	$(MAKE) -C go test

buildout: buildstatic
	$(MAKE) -C out all
