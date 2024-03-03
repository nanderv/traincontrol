#include <EEPROM.h>
#include "globals.h"
#include "coms.h"
#include "startup.h"

void fetchMac()
{
  uint32_t m2 = HW_OCOTP_MAC0;

  MAC[0] = m2 >> 16;
  MAC[1] = m2 >> 8;
  MAC[2] = m2 >> 0;
}
void handleStartState(){
  static messageSlot startupSendSlot;
  fetchMac();

  startupSendSlot.type=0;
  startupSendSlot.content[0]=1;
  startupSendSlot.content[1]=MAC[0];
  startupSendSlot.content[2]=MAC[1];
  startupSendSlot.content[3]=MAC[2];
  setCheckByte(&startupSendSlot);
  writeMessageToAllBut(999, &startupSendSlot);
}

bool handleZeroMode(messageSlot *handleMessage, messageSlot *sendBack){
  if (handleMessage -> content[0] == 254){
    return restart(handleMessage, sendBack);
  }
  return false;
}


bool restart(messageSlot *handleMessage, messageSlot *sendBack) {
  if (handleMessage->content[1] == MAC[0] && handleMessage->content[2] == MAC[1] && handleMessage->content[3] == MAC[2]){
    delay(2000);
    SCB_AIRCR = 0x05FA0004;
  }
  return false;
}
