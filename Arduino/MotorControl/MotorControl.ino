uint8_t pwm[] = {2,3,4,5};
uint8_t dir[] = {22,23,24,25,26,27,28,29};

void setup() {
  Serial.begin(115200);           //Debug (USB)
  Serial2.begin(9600);          //Teensy
  Serial3.begin(300,SERIAL_8O2);  //Tether
  for(int i=0; i<4; i++) {
    analogWrite(pwm[i],0);
  }
  for(int i=0; i<8; i++) {
    digitalWrite(dir[i],LOW); //Everything low
  }
  Serial.println("Started.");
}

void setMotor(int mot, int val) {
  int aval = val;
  if(val<0){
    aval=-val;
  }
  aval*=2;
  if(aval==1) {
    aval = val = 0;
  }
  if(aval>255) {
    aval=255;
  }
  mot--;
  if((mot>3)||(mot<0)) {
    Serial.println("ERROR: invalid motor");
    return;
  }
  int a,b;
  if(val>0) {
    a=HIGH;
    b=LOW;
  }else if(val<0){
    a=LOW;
    b=HIGH;
  }else{
    a=b=LOW;
  }
  digitalWrite(dir[mot*2],a);
  digitalWrite(dir[(mot*2)+1],b);
  analogWrite(pwm[mot],aval);
  Serial.print("Set motor ");
  Serial.print(mot);
  Serial.print(" to ");
  Serial.print(val);
  Serial.println(".");
}

int mot = 0;
void loop() {
  if(mot == 0) {
    mot = Serial2.parseInt();
  }else{
    int val = Serial2.parseInt();
    if(val!=0) {
      setMotor(mot,val);
      mot=0;
    }
  }
  while(Serial3.available()>0) {
    int rd = Serial3.read();
    Serial.println(rd);
    Serial2.write(rd);
  }
}
