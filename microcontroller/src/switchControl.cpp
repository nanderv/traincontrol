#include "switchControl.h"

bool controlSwitch(messageSlot *handleMessage, messageSlot *sendBack) {
  // If the message is not for me
  if (DEVICE_ID != handleMessage->content[0]){
    return false;
  }

  if (handleMessage->content[1] != 0){
    // we do not currently support multiple banks
    return false;
  }

  digitalWrite(handleMessage->content[2], HIGH);
  delay(100);
  digitalWrite(handleMessage->content[2], LOW);


  return sendAck(handleMessage, sendBack, SWITCH_RETURN_TYPE);
}
