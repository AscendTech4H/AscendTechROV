#include <CustomSoftwareSerial.h>

CustomSoftwareSerial* Serial3;

void setup() {
  Serial3 = new CustomSoftwareSerial(10,11);
  Serial.begin(115200);
  Serial3->begin(300, CSERIAL_8O2);
  Serial.println("Init");
}

void loop() {
  int i = Serial.read();
  if(i != -1) {
    Serial3->write(i);
    Serial.print("Write ");
    Serial.println(i);
  }
}
