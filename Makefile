all: buildout

buildstatic:
	$(MAKE) -C static all

clean: cleanstatic cleanout

cleanstatic:
	$(MAKE) -C static clean

cleanout:
	$(MAKE) -C out clean

ci: buildstatic
	$(MAKE) -C out ci
	$(MAKE) clean

buildout: buildstatic
	$(MAKE) -C out all
