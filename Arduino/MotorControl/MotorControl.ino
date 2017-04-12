#include <MsTimer2.h>

uint8_t motp[16];
uint8_t pins[] = {35,34,31,30,52,51,53,50,29,28,23,22,48,49,40,45};

void setup() {
  Serial.begin(115200);
  for(int i=0; i<16; i++) {
    motp[i]=0;
  }
  MsTimer2::set(1,motors); //Run motor update every 1 ms
  MsTimer2::start();
  Serial.println("Started.");
}

uint8_t cnt = 0;
void motors() {
  for(int i=0; i<16; i++) {
    digitalWrite(pins[i],motp[i]>=cnt);
  }
  cnt++;
  Serial.println("motup");
}

void setMotor(int mot, int val) {
  if(mot>7) {
    Serial.println("ERROR: invalid motor");
    return;
  }
  int mval = map(val,0,255,-255,255);
  int a,b;
  if(mval>0) {
    a=mval;
    b=0;
  }else{
    a=0;
    b=-mval;
  }
  motp[mot*2]=a;
  motp[(mot*2)+1]=b;
  Serial.print("Set motor ");
  Serial.print(mot);
  Serial.print(" to ");
  Serial.print(val);
  Serial.println(".");
}

int mot = -1;
void loop() {
  if(mot == -1) {
    mot = Serial.parseInt();
  }else{
    int val = Serial.parseInt();
    if(val!=-1) {
      setMotor(mot,val);
      mot=-1;
    }
  }
}
