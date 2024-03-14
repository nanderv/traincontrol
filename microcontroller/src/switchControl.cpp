#include "switchControl.h"

bool controlSwitch(messageSlot *handleMessage, messageSlot *sendBack) {
  // If the message is not for me
  if (MAC[0] != handleMessage->content[0] || MAC[1] != handleMessage->content[1] || MAC[2] != handleMessage->content[2]){
    return false;
  }

  if (handleMessage->content[3] != 0){
    // we do not currently support multiple banks
    return false;
  }

  digitalWrite(handleMessage->content[4], HIGH);
  delay(100);
  digitalWrite(handleMessage->content[4], LOW);


  return sendAck(handleMessage, sendBack, SWITCH_RETURN_TYPE);
}
