
 void teensyMAC(uint8_t *mac)
{
  uint32_t m1 = HW_OCOTP_MAC1;
  uint32_t m2 = HW_OCOTP_MAC0;
  mac[0] = m1 >> 8;
  mac[1] = m1 >> 0;
  mac[2] = m2 >> 24;
  mac[3] = m2 >> 16;
  mac[4] = m2 >> 8;
  mac[5] = m2 >> 0;
}
void handleNonRunningState(){
  uint8_t mac[6];
  teensyMAC(mac);

  startupSendSlot.type=0;
  startupSendSlot.content[0]=ADDR_BROADCAST;
  startupSendSlot.content[1]=mac[3];
  startupSendSlot.content[2]=mac[4];
  startupSendSlot.content[3]=mac[5];
  startupSendSlot.content[4]=MY_ID;
  startupSendSlot.content[5]=controllerStatus; // 0
  setCheckByte(&startupSendSlot);
  writeMessageToAllBut(999, &startupSendSlot);
}

bool handleZeroMode(messageSlot *handleMessage, messageSlot *sendBack){
  if (handleMessage -> content[0] == ADDR_SET){
    return updateID(handleMessage, sendBack);
  }
  if (handleMessage -> content[0] == RESTART_CODE){
    return restart(handleMessage, sendBack);
  }
  return false;
}
bool updateID(messageSlot *handleMessage, messageSlot *sendBack){
  uint8_t mac[6];
  teensyMAC(mac);
  byte m1 = mac[3];
  byte m2 = mac[4];
  byte m3 = mac[5];
  if (handleMessage->content[1] == m1 && handleMessage->content[2] == m2 && handleMessage->content[3] == m3){
      sendBack -> type = 0;
      sendBack -> content[0] = 1;
      sendBack -> content[1]= m1;
      sendBack -> content[2]= m2;
      sendBack -> content[3]= m3;
      sendBack -> content[4]= handleMessage -> content[4];
      if (MY_ID != handleMessage -> content[4]){
        MY_ID = handleMessage -> content[4];
        EEPROM.write(0, MY_ID);
      }
      controllerStatus = handleMessage -> content[5];
      sendBack -> content[5] = handleMessage -> content[5];
      setCheckByte(sendBack);
      return true;
    }
    return false;
}

bool restart(messageSlot *handleMessage, messageSlot *sendBack){
  if (handleMessage->content[1] == MY_ID){
    delay(2000);
    SCB_AIRCR = 0x05FA0004;
  }
  return false;
}