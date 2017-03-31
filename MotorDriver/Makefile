all: driver.hex
driver.o:
	avr-gcc -Os -mmcu=attiny88 driver.c -o driver.o

driver.hex: driver.o
	avr-objcopy -j .text -j .data -O ihex driver.o driver.hex

clean:
	rm driver.o driver.hex

