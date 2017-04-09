avrdude -P $1 -p attiny88 -c arduino -e -U flash:w:driver.hex -v -b 19200
