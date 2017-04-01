
#include <Servo.h>
#include <Metro.h>
#include <FlexCAN.h>
#include <Wire.h>

 #define PIN0 D, 0
 #define PIN1 D, 2
 #define PIN2 D, 3
 #define PIN3 D, 4
 #define PIN4 B, 7
 #define PIN5 D, 5
 #define PIN6 D, 6
 #define PIN7 B, 0
 #define PIN8 C, 4
 #define PIN9 C, 2
 #define PIN10 C, 1
 #define PIN11 C, 0
 #define PIN12 B, 5
 #define PIN13 B, 4
 #define PIN14 B, 3
 #define PIN15 B, 1

Servo grab;
Metro sysTimer = Metro(1);

int laser = 4, valve = 5;

FlexCAN CANbus(125000); //change baud rate later
static CAN_message_t msg,rxmsg;
static byte input [8];

int txCount;
unsigned int txTimer,rxTimer;

void setup() {
  //CAN BUS SETUP
  CANbus.begin();
  sysTimer.reset();
  
  //LASER SETUP
  pinMode(laser, OUTPUT);
  digitalWrite(laser, LOW);
    
  //SERVO SETUP
  grab.attach(23); //revise pin number later
}

void getMessage() {  
  if (!rxTimer) {
    while (CANbus.read(rxmsg)) {
      for (int i = 0; i < 8; i++) {
        input[i] = rxmsg.buf[i];
      } 
    }
  }
}

void claw(int motor, int value) { //closed = 180, open = 90
    analogWrite(motor, value);    
}

//PROCESS ATTINY INPUT
void processMotor(int motor, int value){ 
  int v = map(value,0,255,-255,255);
  int l,r;
  if(v>0) {
    l = v;
    r = 0;
  } else {
    l = 0;
    r = v;
  }
  
  Wire.beginTransmission(0);
  Wire.write((motor*2));
  Wire.write(l);
  Wire.write((motor*2)+1);
  Wire.write(r);
  Wire.endTransmission();
}

void shineLaser() {  
  digitalWrite(laser, HIGH);
}

void loop(){ 
    getMessage();
    switch (input[0]) {
      case 0: processMotor(input[1], input[2]);
      case 1: shineLaser(); break; 
      case 2: claw(input[1], input[2]); break;

    }    
}  


/*
void sendMessage(byte input [8]) {
  if (!txTimer) {
    txTimer = 100;
    msg.len = 8;
    msg.id = 0x222; //replace with random 

    for (int i = 0; i < 8; i++) {
      msg.buf[i] = input[i]; //input data
    }

    txCount = 6;
    while (txCount-- ) {
      CANbus.write(msg);
      msg.buf[0]++;
    }
  }
}*/
