#include <EEPROM.h>
#include "Arduino.h"
#include "defines.h"
#include "coms.h"
#include "startup.h"
#include "globals.h"


void setup() {
  Coms0.begin(115200);
  Coms1.begin(9600);
  Coms2.begin(9600);
  Coms5.begin(9600);
//   LoadMemory();

  pinMode(LED_BUILTIN, OUTPUT);

  pinMode(2, OUTPUT);
  pinMode(3, OUTPUT);
  pinMode(4, OUTPUT);
  pinMode(5, OUTPUT);

  pinMode(9, OUTPUT);
  pinMode(10, OUTPUT);
  pinMode(11, OUTPUT);
  pinMode(12, OUTPUT);
  comms[0].write = coms0Write;
  comms[0].available = coms0Available;
  comms[0].read = coms0Read;
  comms[0].inUse = true;
  comms[1].write = coms1Write;
  comms[1].available = coms1Available;
  comms[1].read = coms1Read;
  comms[1].inUse = true;
  comms[2].write = coms2Write;
  comms[2].available = coms2Available;
  comms[2].read = coms2Read;
  comms[2].inUse = true;

  comms[3].write = coms5Write;
  comms[3].available = coms5Available;
  comms[3].read = coms5Read;
  comms[3].inUse = true;
}


void loop() {
  static int sleepCount;
  for (int i = 0; i < CHAN_IN_USE; i++) {
    handleChannel(i);
  }
  if(sleepCount ==0){
    digitalWrite(LED_BUILTIN, HIGH);
    handleStartState();
  }

  if(sleepCount == 1000){
    digitalWrite(LED_BUILTIN, LOW);
  }
  sleepCount++;
  if (sleepCount == 5000){
    sleepCount = 0;
  }
  delay(1);
  Serial.flush();
}
