all: out.tar.xz

topSide.o:
	cp ../go/main/main topSide.o

driver.hex:
	cp ../MotorDriver/driver.hex driver.hex

static:
	cp -r ../static static
	rm static/Makefile static/.gitignore

out.tar.xz: topSide.o driver.hex static
	tar -cvf out.tar.xz topSide.o driver.hex static

clean:
	rm -r out.tar.xz topSide.o driver.hex static

makeall:
	$(MAKE) -C ..

cleanall:
	$(MAKE) -C .. clean

rebuildall:
	$(MAKE) cleanall
	$(MAKE) makeall
