#include <inttypes.h>
#include <avr/interrupt.h>
#include <avr/sleep.h>
#include <util/twi.h>

uint8 pinstates[16];
uint8 ctr=0;

//Macro to write to a pin
#define pinWrite(port,pin,val) PORT##port =  val * (val << pin)

//Define rest of pins later
#define PIN0 D, 0
#define PIN1 D, 1
#define PIN2 D, 2
#define PIN3 D, 3
#define PIN4 D, 4
#define PIN5 B, 7
#define PIN6 D, 5
#define PIN7 D, 6
#define PIN8 D, 7
#define PIN9 B, 0
#define PIN10 C, 4
#define PIN11 C, 2
#define PIN12 C, 1
#define PIN13 C, 0
#define PIN14 B, 5
#define PIN15 B, 4
#define PIN16 B, 3
#define PIN17 B, 1

//Pin selection macro (e.g. PINx(0) -> PIN0)
#define PINx(pin) PIN##pin

//Process a pin
#define PINproc(pin) pinWrite( PINx(pin) , pinstates[ pin ] < ctr )

ISR (TIMER1_OVF_vect) { //On timer overflow, process pins
	PINproc(0);
	PINproc(1);
	PINproc(2);
	PINproc(3);
	PINproc(4);
	PINproc(5);
	PINproc(6);
	PINproc(7);
	PINproc(8);
	PINproc(9);
	PINproc(10);
	PINproc(11);
	PINproc(12);
	PINproc(13);
	PINproc(14);
	PINproc(15);
	ctr++;
	TWCR |= b10000000;
}

void setupTimer() {
	TCCR1A = TIMER1_PWM_INIT;
	TCCR1B |= TIMER1_CLOCKSOURCE;
}

void setupTWI() {
	TWAR = 40;		//Set slave address
	TWBR = 8;		//Set speed to 128kbps
	TWCR |= b11000101;	//Set TWI control flags to enable TWI
}

const int ADDR = 0;
const int VAL = 1;
int editstat=ADDR;
int motid;

ISR (TWI_vect) {
	if(editstat==ADDR) {
		motid=TWDR;
		editstat=VAL;
	}else if(editstat==VAL) {
		pinstates[motid]=TWDR;
	}
	TWCR |= b10000000;
}

void main() {
	setupTimer();
	setupTWI();
	while(true){
		sleep_mode();
	}
}
