#include "message.h"
#include "cmds.h"
#include "switchControl.h"
#include "sectorControl.h"

void setCheckByte(messageSlot *m) {
  byte check = m->type;
  for (int i = 0; i < MSG_LENGTH; i++) {
    check = check ^ m->content[i];
  }
  m->checkByte = check;
}

bool handleMessage(messageSlot *in, messageSlot *out) {
  switch(in->type){
    case FLASHING_LIGHTS_TYPE:
      return flashingLights(in, out);
    case SWITCH_TYPE:
      return controlSwitch(in, out);
    case SECTOR_TYPE:
      return controlSector(in, out);
  }
  return false;
}

bool flashingLights(messageSlot *handleMessage, messageSlot *sendBack) {
  digitalWrite(LED_BUILTIN, HIGH);
  delay(handleMessage->content[0] * 10);
  digitalWrite(LED_BUILTIN, LOW);
  sendBack->type = 2;
  for (int i = 0; i < MSG_LENGTH; i++) {
    sendBack->content[i] = handleMessage->content[i];
  }
  setCheckByte(sendBack);
  return false;
}

bool sendAck(messageSlot *handleMessage, messageSlot *sendBack, byte id) {
  sendBack->type = id;
  for (int i = 0; i < MSG_LENGTH; i++) {
    sendBack->content[i] = handleMessage->content[i];
  }
  setCheckByte(sendBack);
  return true;
}
