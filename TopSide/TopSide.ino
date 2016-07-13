#include <SoftwareSerial.h>

const int NUMBER_MOTORS = 8;
const int NUMBER_CAMERAS = 4;
const int NUMBER_SERVOS = 3;

SoftwareSerial Comm(9,11);  //Software serial port for RS-485
const int rsc=10; //RS-485 Control Pin

void readVals(){
  Comm.write(((uint8_t)3)); //Send command 3 (send values)
  delay(20);                //Wait 10ms for bottom side controller to recieve message
  digitalWrite(rsc,LOW);    //Switch to recieve mode
  while(Comm.available()<1){} //Wait for size of sensor values
  uint8_t s=Comm.read();
  while(Comm.available()<s){} //Wait for it to come
  int i;
  for(i=0,i<sizeof(sensors),i++){
    Serial.write((uint8_t)Comm.read()); //Pass it all right on through
  }
  digitalWrite(rsc,HIGH); //Switch back to sending mode
}

void updateBMP(){
  Comm.write(5);
}

void updateExternalTempSensor(){
  Comm.write(6);
}

void updateMotor(){
  while(Serial.available()<2){}
  Comm.write((byte)0);
  Comm.write((byte)Serial.read());
  Comm.write((byte)Serial.read());
}

int cameraPins[NUMBER_CAMERAS] = {4, 5, 6, 7};
void setCamera(int camera){ // turn off all cameras except for the one requested
  for(int i = 0; i < NUMBER_CAMERAS; i++){
    int pin = cameraPins[i];
    if(i == camera){ // camera value will be 0, 1, 2, or 3
      digitalWrite(pin, LOW);
    }else{
      digitalWrite(pin, HIGH);
    }
  }
}

void updateBMP(){
  Comm.write(5);
}

void updateExternalTempSensor(){
  Comm.write(6);
}

void setup(){
  Serial.begin(115200);
  Comm.begin(19200);
  digitalWrite(rsc,HIGH); //Start in sending mode
}

void loop(){
  while(Serial.available()<1){} //Wait for an instruction
  switch(Serial.read()){        //Execute instruction
  case 0: //Motor set
    updateMotor();
    break;
  case 1: //Sensor update
    while(Serial.available()<1){}
    switch(Serial.read()){
    case 0: updateBMP();break;
    case 1: updateExternalTempSensor();break;
    }
    break;
  case 2: //Read in values
    readVals();
    break;
  case 3: //Set camera
    while(Serial.available()<1){}
    setCamera(Serial.read());
    break;
  }
}
