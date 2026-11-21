
#include "sectorControl.h"

bool controlSector(messageSlot *handleMessage, messageSlot *sendBack) {
  // If the message is not for me
  if (DEVICE_ID != handleMessage->content[0]){
    return false;
  }

  if (handleMessage->content[1] != 0){
    // we do not currently support multiple banks
    return false;
  }

  if (handleMessage->content[2] > 0){
    digitalWrite(handleMessage->content[3], HIGH);
  } else {
    digitalWrite(handleMessage->content[3], LOW);
  }

  return sendAck(handleMessage, sendBack, SECTOR_RETURN_TYPE);
}
