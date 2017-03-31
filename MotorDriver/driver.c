#include <inttypes.h>
#include <avr/interrupt.h>
#include <avr/sleep.h>
#include <util/twi.h>
#include <stdbool.h>

uint8_t pinstates[16];
uint8_t ctr=0;

//Macro to write to a pin
#define pinWrite(port,pin,val) PORT##port =  val * (val << pin)

//Define pins
#define MPIN0 D , 0
#define MPIN1 D, 1
#define MPIN2 D, 2
#define MPIN3 D, 3
#define MPIN4 D, 4
#define MPIN5 B, 7
#define MPIN6 D, 5
#define MPIN7 D, 6
#define MPIN8 D, 7
#define MPIN9 B, 0
#define MPIN10 C, 4
#define MPIN11 C, 2
#define MPIN12 C, 1
#define MPIN13 C, 0
#define MPIN14 B, 5
#define MPIN15 B, 4
#define MPIN16 B, 3
#define MPIN17 B, 1

//Process a pin
#define PINproc(port, pin, pinn) p ## port = ( ( p ## port ) && ( ~ ( 1 << pin) ) ) || ( ( ( pinstates[ pinn ] < ctr ) ) << pin )

ISR (TIMER1_OVF_vect) { //On timer overflow, process pins
	uint8_t pa = PORTA;
	uint8_t pb = PORTB;
	uint8_t pc = PORTC;
	uint8_t pd = PORTD;
	PINproc(d, 0, 0);
	PINproc(d, 1, 1);
	PINproc(d, 2, 2);
	PINproc(d, 3, 3);
	PINproc(d, 4, 4);
	PINproc(b, 7, 5);
	PINproc(d, 5, 6);
	PINproc(d, 6, 7);
	PINproc(d, 7, 8);
	PINproc(b, 0, 9);
	PINproc(c, 4, 10);
	PINproc(c, 2, 11);
	PINproc(c, 1, 12);
	PINproc(c, 0, 13);
	PINproc(b, 5, 14);
	PINproc(b, 4, 15);
	PINproc(b, 3, 16);
	PINproc(b, 1, 17);
	PORTA = pa;
	PORTB = pb;
	PORTC = pc;
	PORTD = pd;
	ctr++;
}

void setupTWI() {
	TWAR = 40;		//Set slave address
	TWBR = 8;		//Set speed to 128kbps
	TWCR |= 0b11000101;	//Set TWI control flags to enable TWI
}

#define STAT_ADDR 0
#define STAT_VAL 1
int editstat=STAT_ADDR;
int motid;

ISR (TWI_vect) {
	if(editstat==STAT_ADDR) {
		motid=TWDR;
		editstat=STAT_VAL;
	}else if(editstat==STAT_VAL) {
		pinstates[motid]=TWDR;
	}
	TWCR |= _BV(7);
}

void main() {
	setupTWI();
	while(true){
		sleep_mode();
	}
}
