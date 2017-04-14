all: topSide.o camrelay.o static

topSide.o:
	cp ../go/main/main topSide.o

camrelay.o:
	cp ../go/camrelay/camrelay camrelay.o

static:
	cp -r ../static static
	rm static/Makefile static/.gitignore

clean:
	rm -r topSide.o static camrelay.o

makeall:
	$(MAKE) -C ..

cleanall:
	$(MAKE) -C .. clean

rebuildall:
	$(MAKE) cleanall
	$(MAKE) makeall
