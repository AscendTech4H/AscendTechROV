all: out.tar.xz

topSide.o:
	cp ../go/main/main topSide.o

static:
	cp -r ../static static
	rm static/Makefile static/.gitignore

out.tar.xz: topSide.o static
	tar -cvf out.tar.xz topSide.o static

clean:
	rm -r out.tar.xz static topSide.o
