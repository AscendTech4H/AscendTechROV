#include <Servo.h>

#include <Tlc5940.h>
#include<Wire.h>
#include <OneWire.h>

const int MPU_addr=0x68;  // I2C address of the MPU-6050
const int BMP_ADDRESS=0x77;  // I2C address of the BMP085

const unsigned char oversampling_setting = 3; //oversampling for measurement
const unsigned char pressure_conversiontime[4] = {5,8,14,26};  // delays for oversampling settings 0, 1, 2 and 3
int ac1,ac2,ac3;
unsigned int ac4,ac5,ac6;
int b1,b2;
int mb,mc,md;
int16_t temperature = 0;
long pressure = 0;
float celsius, fahrenheit;

int16_t AcX,AcY,AcZ,Tmp,GyX,GyY,GyZ;

const int switchpin = 16;

OneWire  ds(4);  // on pin 10 (a 4.7K resistor is necessary)
byte addr[8];
int16_t extemp;
byte type_s;

Servo serv0;
Servo serv1;
Servo serv2;

void setup() {
  //DEBUG SETUP
  Serial.begin(115200);
  Serial.println("Init");
  //TETHER SETUP
  Serial3.begin(19200);
  pinMode(16,OUTPUT);
  digitalWrite(16,LOW);
  //TLC SETUP
  Tlc.init();
  Tlc.update();
  //MPU SETUP
  Wire.begin();
  Wire.beginTransmission(MPU_addr);
  Wire.write(0x6B);  // PWR_MGMT_1 register
  Wire.write(0);     // set to zero (wakes up the MPU-6050)
  Wire.endTransmission(true);
  //DMP SETUP
  Wire.begin(0x77);
  getCalibrationData();
  Wire.endTransmission();
  //SERVO SETUP
  serv0.attach(5);
  serv1.attach(6);
  serv2.attach(7);
  //DS18x20 SETUP
  byte i;
  if ( !ds.search(addr)) {
    Serial.println("No more addresses.");
    Serial.println();
    ds.reset_search();
    delay(250);
    return;
  }
  
  Serial.print("ROM =");
  for( i = 0; i < 8; i++) {
    Serial.write(' ');
    Serial.print(addr[i], HEX);
  }

  if (OneWire::crc8(addr, 7) != addr[7]) {
      Serial.println("CRC is not valid!");
      return;
  }
  Serial.println();
 
  // the first ROM byte indicates which chip
  switch (addr[0]) {
    case 0x10:
      Serial.println("  Chip = DS18S20");  // or old DS1820
      type_s = 1;
      break;
    case 0x28:
      Serial.println("  Chip = DS18B20");
      type_s = 0;
      break;
    case 0x22:
      Serial.println("  Chip = DS1822");
      type_s = 0;
      break;
    default:
      Serial.println("Device is not a DS18x20 family device.");
      //return;
  }
  ds.reset();
  //DONE
  Serial.println("Initialization Complete");
}

void  getCalibrationData() {
  Serial.println("Reading Calibration Data");
  ac1 = readIntRegister(0xAA);
  Serial.print("AC1: ");
  Serial.println(ac1,DEC);
  ac2 = readIntRegister(0xAC);
  Serial.print("AC2: ");
  Serial.println(ac2,DEC);
  ac3 = readIntRegister(0xAE);
  Serial.print("AC3: ");
  Serial.println(ac3,DEC);
  ac4 = readIntRegister(0xB0);
  Serial.print("AC4: ");
  Serial.println(ac4,DEC);
  ac5 = readIntRegister(0xB2);
  Serial.print("AC5: ");
  Serial.println(ac5,DEC);
  ac6 = readIntRegister(0xB4);
  Serial.print("AC6: ");
  Serial.println(ac6,DEC);
  b1 = readIntRegister(0xB6);
  Serial.print("B1: ");
  Serial.println(b1,DEC);
  b2 = readIntRegister(0xB8);
  Serial.print("B2: ");
  Serial.println(b1,DEC);
  mb = readIntRegister(0xBA);
  Serial.print("MB: ");
  Serial.println(mb,DEC);
  mc = readIntRegister(0xBC);
  Serial.print("MC: ");
  Serial.println(mc,DEC);
  md = readIntRegister(0xBE);
  Serial.print("MD: ");
  Serial.println(md,DEC);
}

void readBMPINT(){
  int  ut= readUT();
  long up = readUP();
  long x1, x2, x3, b3, b5, b6, p;
  unsigned long b4, b7;

  //calculate true temperature
  x1 = ((long)ut - ac6) * ac5 >> 15;
  x2 = ((long) mc << 11) / (x1 + md);
  b5 = x1 + x2;
  temperature = (b5 + 8) >> 4;

  //calculate true pressure
  b6 = b5 - 4000;
  x1 = (b2 * (b6 * b6 >> 12)) >> 11; 
  x2 = ac2 * b6 >> 11;
  x3 = x1 + x2;
  b3 = (((int32_t) ac1 * 4 + x3)<<oversampling_setting + 2) >> 2;
  x1 = ac3 * b6 >> 13;
  x2 = (b1 * (b6 * b6 >> 12)) >> 16;
  x3 = ((x1 + x2) + 2) >> 2;
  b4 = (ac4 * (uint32_t) (x3 + 32768)) >> 15;
  b7 = ((uint32_t) up - b3) * (50000 >> oversampling_setting);
  p = b7 < 0x80000000 ? (b7 * 2) / b4 : (b7 / b4) * 2;

  x1 = (p >> 8) * (p >> 8);
  x1 = (x1 * 3038) >> 16;
  x2 = (-7357 * p) >> 16;
  pressure = p + ((x1 + x2 + 3791) >> 4);
}

// read uncompensated temperature value
unsigned int readUT() {
  writeRegister(0xf4,0x2e);
  delay(5); // the datasheet suggests 4.5 ms
  return readIntRegister(0xf6);
}

// read uncompensated pressure value
long readUP() {
  writeRegister(0xf4,0x34+(oversampling_setting<<6));
  delay(pressure_conversiontime[oversampling_setting]);

  unsigned char msb, lsb, xlsb;
  Wire.beginTransmission(BMP_ADDRESS);
  Wire.write(0xf6);  // register to read
  Wire.endTransmission();

  Wire.requestFrom(BMP_ADDRESS, 3); // request three bytes
  while(!Wire.available()); // wait until data available
  msb = Wire.read();
  while(!Wire.available()); // wait until data available
  lsb |= Wire.read();
  while(!Wire.available()); // wait until data available
  xlsb |= Wire.read();
  return (((long)msb<<16) | ((long)lsb<<8) | ((long)xlsb)) >>(8-oversampling_setting);
}

void writeRegister(unsigned char r, unsigned char v)
{
  Wire.beginTransmission(BMP_ADDRESS);
  Wire.write(r);
  Wire.write(v);
  Wire.endTransmission();
}

// read a 16 bit register
int readIntRegister(unsigned char r)
{
  unsigned char msb, lsb;
  Wire.beginTransmission(BMP_ADDRESS);
  Wire.write(r);  // register to read
  Wire.endTransmission();

  Wire.requestFrom(BMP_ADDRESS, 2); // request two bytes
  while(!Wire.available()); // wait until data available
  msb = Wire.read();
  while(!Wire.available()); // wait until data available
  lsb = Wire.read();
  return (((int)msb<<8) | ((int)lsb));
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
  Tlc.set(motor,a);
  Tlc.set(motor+1,b);
}

void processMPURead(){
  Serial.println("MPU READ");
  Wire.beginTransmission(MPU_addr);
  Wire.write(0x3B);  // starting with register 0x3B (ACCEL_XOUT_H)
  Wire.endTransmission(false);
  Wire.requestFrom(MPU_addr,14,true);  // request a total of 14 registers
  AcX=Wire.read()<<8|Wire.read();  // 0x3B (ACCEL_XOUT_H) & 0x3C (ACCEL_XOUT_L)    
  AcY=Wire.read()<<8|Wire.read();  // 0x3D (ACCEL_YOUT_H) & 0x3E (ACCEL_YOUT_L)
  AcZ=Wire.read()<<8|Wire.read();  // 0x3F (ACCEL_ZOUT_H) & 0x40 (ACCEL_ZOUT_L)
  Tmp=Wire.read()<<8|Wire.read();  // 0x41 (TEMP_OUT_H) & 0x42 (TEMP_OUT_L)
  GyX=Wire.read()<<8|Wire.read();  // 0x43 (GYRO_XOUT_H) & 0x44 (GYRO_XOUT_L)
  GyY=Wire.read()<<8|Wire.read();  // 0x45 (GYRO_YOUT_H) & 0x46 (GYRO_YOUT_L)
  GyZ=Wire.read()<<8|Wire.read();  // 0x47 (GYRO_ZOUT_H) & 0x48 (GYRO_ZOUT_L)
  Serial.print("AcX = "); Serial.print(AcX);
  Serial.print(" | AcY = "); Serial.print(AcY);
  Serial.print(" | AcZ = "); Serial.print(AcZ);
  Serial.print(" | Tmp = "); Serial.print(Tmp/340.00+36.53);  //equation for temperature in degrees C from datasheet
  Serial.print(" | GyX = "); Serial.print(GyX);
  Serial.print(" | GyY = "); Serial.print(GyY);
  Serial.print(" | GyZ = "); Serial.println(GyZ);
}

void sendData(){
  Serial.println("Sending Data");
  digitalWrite(16,HIGH);
  int16_t data[9]={AcX,AcY,AcZ,Tmp,GyX,GyY,GyZ,temperature,extemp};
  byte* databytes=(byte*)(&data);
  int i;
  for(i=0;i<18;i++){
    Serial3.write(databytes[i]);
  }
  databytes=(byte*)(&pressure);
  for(i=0;i<4;i++){
    Serial3.write(databytes[i]);
  }
  delay(100);
  digitalWrite(16,LOW);
}

void bmpProcess(){
  readBMPINT();
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
    case 1: serv1.write(servValue); break;
    case 2: serv2.write(servValue); break;
  }
}

void extempProcess(){
  byte present = 0;
  byte data[12];
  byte i;
  ds.reset();
  ds.select(addr);
  ds.write(0x44);
  present = ds.reset();
  ds.select(addr);    
  ds.write(0xBE);         // Read Scratchpad
  for ( i = 0; i < 9; i++) {           // we need 9 bytes
    data[i] = ds.read();
    //Serial.print(data[i], HEX);
    //Serial.print(" ");
  }
  //Serial.print(" CRC=");
  //Serial.print(OneWire::crc8(data, 8), HEX);
  //Serial.println();
  int16_t raw = (data[1] << 8) | data[0];
  if (type_s) {
    raw = raw << 3; // 9 bit resolution default
    if (data[7] == 0x10) {
      // "count remain" gives full 12 bit resolution
      raw = (raw & 0xFFF0) + 12 - data[6];
    }
  } else {
    byte cfg = (data[4] & 0x60);
    // at lower res, the low bits are undefined, so let's zero them
    if (cfg == 0x00) raw = raw & ~7;  // 9 bit resolution, 93.75 ms
    else if (cfg == 0x20) raw = raw & ~3; // 10 bit res, 187.5 ms
    else if (cfg == 0x40) raw = raw & ~1; // 11 bit res, 375 ms
    //// default is 12 bit resolution, 750 ms conversion time
  }
  extemp=raw;
}

void tetherProcess(){
  Serial.println("Processing Command");
  //Serial.println(Serial3.read());
  switch(Serial3.read()){ //PROCESS COMMAND
    case 0: processMotor();break;                               //TLC INPUT
    case 1: Serial.println("TLC Update");Tlc.update();break;    //TLC UPDATE
    case 2: processMPURead();break;                             //MPU READ
    case 3: sendData();break;                                   //SEND DATA
    case 4: servoSet();break;                                   //SET SERVO INPUTS
    case 5: bmpProcess();break;                                 //BMP READ
    case 6: Serial.println("Update Extemp");extempProcess();break;
  }
}

void loop(){
  if(Serial3.available()>0){tetherProcess();}
}
