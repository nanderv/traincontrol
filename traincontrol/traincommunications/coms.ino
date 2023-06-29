#include "hexConv.h"



void writeToAllBut(int dest, char input) {
  for (int i = 0; i < CHAN_IN_USE; i++) {
    if (comss[i].inUse && (DEBUG_SENDBACK || i != dest)) {
      twoHex crs = fromHex(input);
      comss[i].write(crs.fst);
      comss[i].write(crs.snd);
    }
  }
}

void writeMessageToAllBut(int dest, messageSlot *msg ) {
  char checkByte;
  checkByte = 0;
  writeToAllBut(dest, msg->type);
  checkByte = checkByte ^ msg-> type;
  for (int i = 0; i < MSG_LENGTH; i++) {
    char bt = msg->content[i];
    writeToAllBut(dest, bt);
    checkByte = checkByte ^ bt;
  }
  writeToAllBut(dest, checkByte);
}


void resetReader(int dest) {
  // Reset the half read
  comss[dest].halfRead = 0;
  comss[dest].isHalfRead = false;
  // And reset the message that is read
  comss[dest].msg.readState = type;
}


// True if finished, otherwise false.
bool readMessage(int comNR) {
  messageSlot *slot = &(comss[comNR].msg);
  byte checkByte = 0;
  if (!toHexDuo(comss[comNR].read(), comss[comNR].read(), &slot->type)){
    return false;
  }
  checkByte = slot ->type;

  for(int i = 0; i< MSG_LENGTH;i++){
      if (!toHexDuo(comss[comNR].read(), comss[comNR].read(), &slot->content[i])){
        return false;
      };
      checkByte = checkByte ^ slot->content[i];
  }

  if (!toHexDuo(comss[comNR].read(), comss[comNR].read(), &slot->checkByte)){
    return false;
  }
  return checkByte == slot->checkByte;
}

void handleChannel(int i){
   if (comss[i].inUse) {
      if(comss[i].available()>=(2*(MSG_LENGTH+2))) {
        while (comss[i].available()>=(2*(MSG_LENGTH+2))) {
          if (readMessage(i)){
            writeMessageToAllBut( 4, &(comss[i].msg));
          } else {
            writeMessageToAllBut( 4, &(comss[i].msg));
          }
        }
        while (comss[i].available()){
          comss[i].read();
        }
      }
   }
}


byte coms1Read() {
  return Coms1.read();
}

int coms1Available() {
  return Coms1.available();
}

void coms1Write(char c) {
  Coms1.write(c);
}

char voidRead(){return '-';}
bool voidAvailable() {return false;}
void voidWrite(char c) {};
