# define Coms0 Serial
# define Coms1 Serial1
# define Coms3 Serial3
# define Coms4 Serial4
# define Coms5 Serial5
# define CHAN_IN_USE 6
# define ERROR_TYPE 255
#include "cmds.h"
#include "coms.h"
#include "metaControl.h"
#include "switchControl.h"
#include "sectorControl.h"
#include "startup.h"

void setup() {
  addHandler(0,metaControl );
  addHandler(1, flashingLights);
  addHandler(SWITCH_TYPE, controlSwitch);
  addHandler(SECTOR_TYPE, controlSector);

  Coms0.begin(115200);
  Coms1.begin(9600);
  Coms3.begin(9600);
  Coms4.begin(9600);
  Coms5.begin(9600);
  LoadMemory();

  pinMode(LED_BUILTIN, OUTPUT);
  
  comms[0].write = coms0Write;
  comms[0].available = coms0Available;
  comms[0].read = coms0Read;
  comms[0].inUse = true;
  comms[1].write = coms1Write;
  comms[1].available = coms1Available;
  comms[1].read = coms1Read;
  comms[1].inUse = true;
  comms[2].write = coms3Write;
  comms[2].available = coms3Available;
  comms[2].read = coms3Read;
  comms[2].inUse = true;
  comms[3].write = coms4Write;
  comms[3].available = coms4Available;
  comms[3].read = coms4Read;
  comms[3].inUse = true;
  comms[4].write = coms5Write;
  comms[4].available = coms5Available;
  comms[4].read = coms5Read;
  comms[4].inUse = true;

}

int controllerStatus ;
int sleepCount;
messageSlot startupSendSlot;

void loop() {
  for (int i = 0; i < CHAN_IN_USE; i++) {
    handleChannel(i);
  }
  if (controllerStatus != STATUS_RUNNING){
    if(sleepCount ==0){
      digitalWrite(LED_BUILTIN, HIGH);
      handleNonRunningState();
    }
    if(sleepCount == 1000){
      digitalWrite(LED_BUILTIN, LOW);
    }
    sleepCount ++;
    if (sleepCount == 5000){
      sleepCount = 0;
    }
    delay(1);
  }
  Serial.flush();
}
