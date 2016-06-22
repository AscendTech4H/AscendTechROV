#include <SoftwareSerial.h>

//#define DEBUG
/*IMPORTANT NOTE ABOUT DEBUG:
  When in DEBUG mode, do not run the python at the same time. The 
  python recieves messages from the Arduino via serial and having
  the Arduino print out data from serial will interfere with the 
  python. Instead, commands must be manually typed into the serial 
  monitor.
  
  To view all sensor values currently stored, type in '&'*/

/*RECEIVING COMMANDS (from Python):
    Motors:
      ^value0, value1, value2, value3...$
    Servo Motors:
      ^value0>value2>value3*
 */

/* commands to bottomside:
  0, motorName, speed - write motor
  1 - tlc update
  2 - read mpu
  3 - send mpu & temp
  4 - servoName, position - write servo
  5 - update bmp
  6 - update external temperature sensor
*/

const int NUMBER_MOTORS = 8;
const int NUMBER_CAMERAS = 4;
const int NUMBER_SERVOS = 3;

//Pin constants
const int TX = 9;
const int RX = 11;
const int TX_CONTROL = 10;

//RS485 constants
const int TRANSMIT = HIGH;
const int RECEIVE = LOW;

SoftwareSerial Comm(TX, RX); // allows for communication with bottomside

struct SensorVals { // struct that encapsulates sensor values
  float tempf;
  float tempc;
  float distance;
  int16_t AcX,AcY,AcZ,intemp,GyX,GyY,GyZ;
  uint16_t pres;
} sensors;

SensorVals sense = sensors; // stores all sensor readings

int cameraPins[NUMBER_CAMERAS] = {4, 5, 6, 7};

int motorVals[NUMBER_MOTORS];
int prevMotorVals[NUMBER_MOTORS]; // used to check which motors have been updated after input is polled

int servoVals[NUMBER_SERVOS];
int prevServoVals[NUMBER_SERVOS];

String current = "";
boolean isReading = false;
int currentMotor;

int16_t AcX,AcY,AcZ,Tmp,GyX,GyY,GyZ,temp2;
long pressure;

void setup()
{
  pinMode(TX_CONTROL, OUTPUT);
  digitalWrite(TX_CONTROL, TRANSMIT);
  
  for (int i = 0; i < NUMBER_CAMERAS; i++){
    pinMode(cameraPins[i], OUTPUT);
  }

  setCamera(0);

  Serial.begin(19200); // will be used to communicate to python code
  Comm.begin(19200); // will be used to communicate with bottomside
  Serial.setTimeout(6400000);
}

void loop()
{
  pollInput();
}

void pollInput() // poll and parse input
{
  while(Serial.available())
  {
    char inChar = (char) Serial.read();
    if(inChar == '^')
    {
      //new command just started
      isReading = true;
      currentMotor = 0;
    }
    else if(inChar == '$')
    {
      //command ended, set the motors
      motorVals[currentMotor] = current.toInt();
      current.remove(0);
      isReading = false;
      sendMotors();
      currentMotor = 0;
    }
    else if(inChar == ',')
    {
      //we finished getting the current motor, parse the value
      motorVals[currentMotor] = current.toInt();
      current.remove(0);
      currentMotor++;
    }
    else if(inChar == 'c')
    {
      //change the camera
      setCamera(Serial.parseInt());
    }
    else if(inChar == 'p')
    {
      requestData();
    }
    else if(inChar == 'b')
    {
      updateBMP();
    }
    else if(inChar == 'e')
    {
      updateExternalTempSensor();
    }
    else if(inChar == '&')
    {
      #ifdef DEBUG
        printDebugData(); //only call if in DEBUG mode
      #endif
    }
    else if(inChar == '>')
    {
      //we finished getting the current servo, parse the value
      servoVals[currentMotor] = current.toInt();
      current.remove(0);
      currentMotor++;
    }
    else if(inChar == '*')
    {
      //command ended, set the motors
      servoVals[currentMotor] = current.toInt();
      current.remove(0);
      currentMotor = 0;
      sendServos();
      isReading = false;
    }
    else if(isReading)
    {
      current += inChar;
    }
  }
}

void writeMotor(byte motor,int mspeed) { // send motor data to bottomside
  #ifdef DEBUG
  Serial.print("\tMotor: ");
  Serial.print(motor);
  Serial.print(", Speed: ");
  Serial.println(mspeed);
  #endif
  
  Comm.write((byte)0);
  Comm.write(motor);
  Comm.write(mspeed+128);
}

void writeServo(byte servo, byte pos) { // send servo data to bottomside
  #ifdef DEBUG
  Serial.print("\tServo: ");
  Serial.print(servo);
  Serial.print(", Position: ");
  Serial.println(pos);
  #endif
  
  Comm.write(4);
  Comm.write(servo);
  Comm.write(pos);
}

void requestData() // ask bottomside to send sensor data
{
  Comm.write(2);
  Comm.write(3); // send it up
  digitalWrite(TX_CONTROL,RECEIVE); // go into receive mode so that bottomside can transmit data
  while (Comm.available() < 22);
  int16_t data[9];
  byte* databytes=(byte*)(&data);
  int i;
  for(i=0;i<18;i++){
    databytes[i]=byte(Comm.read());
  }
  databytes=(byte*)(&pressure);
  for(i=0;i<4;i++){
    databytes[i]=byte(Comm.read());
  }
  digitalWrite(TX_CONTROL, TRANSMIT); // return transmit mode
  
  /*Relay information to the python code to interpret
    (both the size of the value and the actual value will be relayed)*/
  #ifndef DEBUG
  Serial.print("d");
  Serial.print(data[0]); Serial.print(","); //AcX
  Serial.print(data[1]); Serial.print(","); //AcY
  Serial.print(data[2]); Serial.print(","); //AcZ
  Serial.print(data[3]); Serial.print(","); //Tmp
  Serial.print(data[4]); Serial.print(","); //GyX
  Serial.print(data[5]); Serial.print(","); //GyY
  Serial.print(data[6]); Serial.print(","); //GyZ
  //Serial.print(data[7]); Serial.print(","); //other Tmp
  Serial.print(pressure); Serial.print(",");
  Serial.print(data[8]); Serial.print("D");
  #endif

 
}

void updateBMP()
{
  Comm.write(5);
}

void updateExternalTempSensor()
{
  Comm.write(6);
}


void sendMotors() // send updated motor values to bottomside
{
  #ifdef DEBUG
  Serial.println("SENDING MOTORS:");
  #endif
  for(int i = 0; i < NUMBER_MOTORS; i++)
  {
    if(motorVals[i] != prevMotorVals[i]) // only send motor values that have changed
    {
      writeMotor(byte(i), motorVals[i]);
    }
  }

  Comm.write(1);
  updatePrevMotors();
}

void sendServos()
{
  #ifdef DEBUG
  Serial.println("SENDING SERVOS:");
  #endif
  for(byte i = 0; i < NUMBER_SERVOS; i++)
  {
    if(servoVals[i] != prevServoVals[i])
    {
      writeServo(i, byte(servoVals[i]));
    }
  }

  updatePrevServos();
}

void updatePrevMotors() // update prevMotorVals to reflect changes to motorVals
{
  for(int i = 0; i < NUMBER_MOTORS; i++)
  {
    prevMotorVals[i] = motorVals[i];
  }
}

void updatePrevServos() // update prevServoVals to reflect changes to servoVals
{
  for(int i = 0; i < NUMBER_SERVOS; i++)
  {
    prevServoVals[i] = servoVals[i]; 
  }
}

void setCamera(int camera) // turn off all cameras except for the one requested
{
  for(int i = 0; i < NUMBER_CAMERAS; i++)
  {
    int pin = cameraPins[i];
    if(i == camera) // camera value obtained from python will be 0, 1, 2, or 3
    {
      #ifdef DEBUG
      Serial.print("TURNING ON CAMERA #");
      Serial.print(i);
      Serial.print(" @ PIN #");
      Serial.println(pin);
      #endif
      digitalWrite(pin, LOW);
    }
    else
    {
      digitalWrite(pin, HIGH);
    }
  }
}

void printDebugData()
{
  Serial.print("Temp F: ");
  Serial.println(sense.tempf);
  Serial.print("Temp C: ");
  Serial.println(sense.tempc);
  Serial.print("Distance: ");
  Serial.println(sense.distance);
  Serial.print("Accelerometer");
  Serial.print("\tx: ");
  Serial.println(sense.AcX);
  Serial.print("\ty: ");
  Serial.println(sense.AcY);
  Serial.print("\tz: ");
  Serial.println(sense.AcZ);
  Serial.print("Internal Temperature: ");
  Serial.println(sense.intemp);
  Serial.println("Gyroscope");
  Serial.print("\tx: ");
  Serial.println(sense.GyX);
  Serial.print("\ty: ");
  Serial.println(sense.GyY);
  Serial.print("\tz: ");
  Serial.println(sense.GyZ);
  Serial.print("Pressure: ");
  Serial.println(sense.pres);
}
