#include "coms.h"


void setup() {
  // put your setup code here, to run once:
  Serial.begin(115200);
  pinMode(LED_BUILTIN, OUTPUT);
  comss[0].bitmask = 0;
  comss[0].write = coms1Write;
  comss[0].available = coms1Available;
  comss[0].read = coms1Read;
  comss[0].inUse = true;
}


void loop() {
  for (int i = 0; i < CHAN_IN_USE; i++) {
    handleChannel(i);
    }
    Serial.flush();
}
