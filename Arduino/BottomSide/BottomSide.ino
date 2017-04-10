#include <SoftwareSerial.h>
#include <Servo.h>
#include <Metro.h>
#include <FlexCAN.h>
#include <Wire.h>

Servo grab;
Metro sysTimer = Metro(1);

int laser = 4, valve = 5;

SoftwareSerial motorControl(7,8);
FlexCAN CANbus(125000); //change baud rate later
static CAN_message_t rxmsg;
static byte input [8];

int txCount;
unsigned int txTimer,rxTimer;

void setup() {
  //CAN BUS SETUP
  CANbus.begin();
  sysTimer.reset();

  //Moror control setup
  motorControl.begin(115200);
  
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

//Update motor states
void processMotor(int motor, int value){
  if(motor<8) {
    motorControl.println(motor);
    motorControl.println(value);
  }else if(motor==9) {
    int val = map(value,0,255,-255,255);
    int l, r;
    if(val>0){
      l = val;
      r = 0;
    }else{
      l = 0;
      r = -val;
    }
    analogWrite(5,l);
    analogWrite(6,r);
  }else if(motor==10){
    grab.write(value);
  }
}

void shineLaser(int state) {  
  digitalWrite(laser, state);
}

void loop(){ 
    getMessage();
    switch (input[0]) {
      case 0: processMotor(input[1], input[2]);
      case 1: shineLaser(input[1]); break;
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
