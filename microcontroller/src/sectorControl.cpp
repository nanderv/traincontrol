
#include "sectorControl.h"

bool controlSector(messageSlot *handleMessage, messageSlot *sendBack) {
  // If the message is not for me
  if (MAC[0] != handleMessage->content[0] && MAC[1] != handleMessage->content[1] && MAC[2] != handleMessage->content[2]){
    return false;
  }

  if (handleMessage->content[3] != 0){
    // we do not currently support multiple banks
    return false;
  }

  if (handleMessage->content[5] > 0){
    digitalWrite(handleMessage->content[4], HIGH);
  } else {
    digitalWrite(handleMessage->content[4], LOW);
  }

  return sendAck(handleMessage, sendBack, SECTOR_RETURN_TYPE);
}
