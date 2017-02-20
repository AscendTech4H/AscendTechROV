#include <Servo.h>
//#include <Tlc5940.h>

Servo serv0;
const int laser = 4;
const int claw = 5;

void setup() {
  //DEBUG SETUP
  Serial.begin(115200);
  Serial.println("Init");
  
  //TETHER SETUP
  Serial3.begin(19200);
  
  pinMode(16,OUTPUT);
  digitalWrite(16,LOW);

  pinMode(laser, OUTPUT);
  digitalWrite(laser, LOW);
  
  /*TLC SETUP
    Tlc.init();
    Tlc.update(); */
    
  //SERVO SETUP
  serv0.attach(23);
  
  Serial.println("Initialization Complete");
}

void turnValve() {
  Serial3.println("Turning valve");
  while (Serial3.available() < 2);
  byte valve = Serial3.read();

  switch(valve) {
    case 7: analogWrite(claw, 128);
            delay(3000);
  }
}


//PROCESS TLC INPUT
void processMotor(){
  Serial.println("TLC Set");
  while (Serial3.available() < 2);
  byte motor = 2*Serial3.read();
  char value = map(Serial3.read(),0,255,-128,127);
  Serial.print("Motor ");
  Serial.print(motor/2);
  Serial.print(" = ");
  Serial.println(value);
  int a, b;
  if (value < 0) {
    a = map(-long(value), 0, 128, 0, 4095);
    b = 0;
  } else {
    a = 0;
    b = map(long(value), 0, 127, 0, 4095);
  }
//  Tlc.set(motor,a);
//  Tlc.set(motor+1,b);
}

void shineLaser() {
  
  while(Serial3.available()<2);
  byte shine = Serial3.read();

  switch(shine) {
    case 6: digitalWrite(laser, HIGH);
  }  
}
void servoSet(){
  Serial.print("Set servo ");
  while(Serial3.available()<2);
  byte servIndex = Serial3.read();
  Serial.print(servIndex);
  byte servValue = map(Serial3.read(), 0, 255, 0, 180);
  Serial.print(" to ");
  Serial.println(servValue);
  switch (servIndex){
    case 0: serv0.write(servValue); break;
  }
}

void tetherProcess(){
  Serial.println("Processing Command");
  //Serial.println(Serial3.read());
  switch(Serial3.read()){		//PROCESS COMMAND
    case 0: processMotor();break;				//TLC INPUT
    case 1: Serial.println("TLCUpdate");//Tlc.update();break;	//TLC UPDATE
    case 4: servoSet();break;					//SET SERVO INPUTS
  }
}

void loop(){
  if(Serial3.available()>0){tetherProcess();}
}
