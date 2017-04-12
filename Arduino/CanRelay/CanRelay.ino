#include <mcp_can.h>
#include <mcp_can_dfs.h>

MCP_CAN CAN0(24);

void setup() {
  Serial.begin(115200);
  Serial.print("Init: ");
  Serial.println(CAN0.begin(MCP_ANY, CAN_125KBPS, MCP_8MHZ) == CAN_OK); //Print success or faliure
  CAN0.setMode(MCP_NORMAL);
}

void loop() {
  int len=0;
  while(len == 0) {
    len = Serial.parseInt();
  }
  byte msg[len];
  for(int n=0; n<len;) {
    int scan = Serial.parseInt();
    if(scan != 0) {
      msg[n]=scan-1;
      n++;
    }
  }
  Serial.print("Send: ");
  Serial.println(CAN0.sendMsgBuf(0,len,(byte*)&msg));
}
