#define SECTOR_EEPROM_MODE 2

bool controlSector(messageSlot *handleMessage, messageSlot *sendBack){
  int ptr = 0;
  bool foundOne=false;
  for (;ptr<MAX_SLOTS; ptr++){
    // IF it's a switch
    // AND the switch has the same ID as the 
      if(controls[ptr].type ==SECTOR_EEPROM_MODE && controls[ptr].data[0] ==  handleMessage->content[0]  ){
        foundOne = true;
        break;
      }
  }
  if (!foundOne){
    return false;
  }

    if (handleMessage->content[1] == 0){
    digitalWrite(controls[ptr].data[1], HIGH);
  } else if (handleMessage->content[1] == 1) {
    digitalWrite(controls[ptr].data[1], LOW);
  } else {
      sendBack->type =ERROR_TYPE;
      sendBack->content[0]=handleMessage->type;
      sendBack->content[1]=handleMessage->content[0];
      sendBack->content[2]=handleMessage->content[1];
      setCheckByte(sendBack);

      return true;
  }
  
  return sendAck(handleMessage, sendBack, SECTOR_RETURN_TYPE);
}
