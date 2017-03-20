all: buildstatic

buildstatic: 
	$(MAKE) -C static

clean: 
	$(MAKE) -C static clean
