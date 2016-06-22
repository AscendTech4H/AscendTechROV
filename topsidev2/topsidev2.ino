#include <SoftwareSerial.h>

SoftwareSerial comm(9, 11);

void setup(){
  pinMode(10,OUTPUT);
  digitalWrite(10,HIGH);
  Serial.begin(115200);
  comm.begin(19200);
}
byte motor=0;
int16_t AcX,AcY,AcZ,Tmp,GyX,GyY,GyZ,temp2;
long pressure;
void loop() {
  Serial.println("Sending Commands");
  comm.write((byte)0);comm.write(motor);comm.write(120);
  comm.write(1);
  comm.write(2);
  comm.write(3);
  delay(1);
  digitalWrite(10,LOW);
  Serial.println("Waiting for response");
  while(comm.available()<20);
  int16_t data[8];
  byte* databytes=(byte*)(&data);
  int i;
  for(i=0;i<16;i++){
    databytes[i]=byte(comm.read());
  }
  databytes=(byte*)(&pressure);
  for(i=0;i<4;i++){
    databytes[i]=byte(comm.read());
  }
  AcX=data[0];
  AcY=data[1];
  AcZ=data[2];
  Tmp=data[3];
  GyX=data[4];
  GyY=data[5];
  GyZ=data[6];
  temp2=data[7];
  //AcX=comm.read()<<8|comm.read();
  //AcY=comm.read()<<8|comm.read();
  //AcZ=comm.read()<<8|comm.read();
  //Tmp=comm.read()<<8|comm.read();
  //GyX=comm.read()<<8|comm.read();
  //GyY=comm.read()<<8|comm.read();
  //GyZ=comm.read()<<8|comm.read();
  Serial.print("AcX = "); Serial.print(AcX);
  Serial.print(" | AcY = "); Serial.print(AcY);
  Serial.print(" | AcZ = "); Serial.print(AcZ);
  Serial.print(" | Tmp = "); Serial.print(Tmp/340.00+36.53);  //equation for temperature in degrees C from datasheet
  Serial.print(" | GyX = "); Serial.print(GyX);
  Serial.print(" | GyY = "); Serial.print(GyY);
  Serial.print(" | GyZ = "); Serial.println(GyZ);
  digitalWrite(10,HIGH);
  delay(2000);
  comm.write((byte)0);comm.write(motor);comm.write((byte)128);
  comm.write(1);
  delay(2000);
  motor++;
  if(motor>7){motor=0;}
}
