#include <SoftwareSerial.h>
#include <Servo.h>

Servo grab;

int laser = 0, valve = 5;

SoftwareSerial motorControl(7,8);

int txCount;
unsigned int txTimer,rxTimer;

void setup() {
  //Moror control setup
  motorControl.begin(9600);
  
  //LASER SETUP
  pinMode(laser, OUTPUT);
  digitalWrite(laser, LOW);

  //SERVO SETUP
  grab.attach(23); //revise pin number later
}

//Update motor states
void processMotor(int motor, int value){
  if(motor<4) {
    motorControl.println(motor+1);
    motorControl.println(value-128);
  }else if(motor==4) {
    value-=128;
    value*=2;
    int aval = value;
    bool neg = (aval<0);
    if(neg) {
      aval*=-1;
    }
    int l, r;
    if(!neg){
      l = aval;
      r = 0;
    }else{
      l = 0;
      r = aval;
    }
    analogWrite(5,l);
    analogWrite(6,r);
  }else if(motor==5){
    grab.write(value);
  }
}

void shineLaser(int state) {  
  digitalWrite(laser, state);
}

void loop(){
  while(motorControl.available()<1){}
  Serial.println("S");
  int rd = motorControl.read();
  Serial.println(rd);
  switch (rd) {
    case 0: while(motorControl.available()<2){}; processMotor(motorControl.read(), motorControl.read());Serial.println("A");break;
    case 1: while(motorControl.available()<1){}; shineLaser(motorControl.read());Serial.println("B");break;
  }
}
