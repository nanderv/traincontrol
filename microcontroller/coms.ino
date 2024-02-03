#include "hexConv.h"



void writeToAllBut(int dest, char *input) {
  for (int i = 0; i < CHAN_IN_USE; i++) {
    if (comms[i].inUse && (DEBUG_SENDBACK || i != dest)) { 
      comms[i].write(input,2*(MSG_LENGTH+2));
    }
  }
}

char messageBuf[2*(MSG_LENGTH+2)];

void convertMessageToByte(messageSlot *msg,  char *buffer) {
  twoHex inB;
  inB = fromHex(msg->type);
  buffer[0] = inB.fst;
  buffer[1] = inB.snd;
  
  for (int i = 0; i < MSG_LENGTH; i++) {
    inB = fromHex(msg->content[i]);
      buffer[2*(i+1)] = inB.fst;
      buffer[2*(i+1)+1] = inB.snd;
  }
  inB = fromHex(msg->checkByte);
  buffer[2*(MSG_LENGTH+1)] = inB.fst;
  buffer[2*(MSG_LENGTH+1)+1] = inB.snd;
}

void writeMessageToAllBut(int dest, messageSlot *msg) {
    convertMessageToByte(msg, messageBuf);
    writeToAllBut(dest, messageBuf);
}


// True if finished, otherwise false.
bool readMessage(int comNR) {
  messageSlot *slot = &(comms[comNR].incomingMessage);
  byte checkByte = 0;
  if (!toHexDuo(comms[comNR].read(), comms[comNR].read(), &slot->type)){
    return false;
  }
  checkByte = slot ->type;

  for(int i = 0; i< MSG_LENGTH;i++){
      if (!toHexDuo(comms[comNR].read(), comms[comNR].read(), &slot->content[i])){
        return false;
      };
      checkByte = checkByte ^ slot->content[i];
  }

  if (!toHexDuo(comms[comNR].read(), comms[comNR].read(), &slot->checkByte)){
    return false;
  }

  if(checkByte == slot->checkByte){
    digitalWrite(LED_BUILTIN, LOW);
  } else {
    digitalWrite(LED_BUILTIN, HIGH);

  }
  return checkByte == slot->checkByte;
}

void handleChannel(int i){
   if (comms[i].inUse) {
      if(comms[i].available()>=(2*(MSG_LENGTH+2))) {
        while (comms[i].available()>=(2*(MSG_LENGTH+2))) {
          if (readMessage(i)){
              writeMessageToAllBut( i, &(comms[i].incomingMessage));
            if (comms[i].incomingMessage.type == 0){
              if(handleZeroMode(&comms[i].incomingMessage, &comms[i].outgoingMessage)){
                writeMessageToAllBut(999, &comms[i].outgoingMessage);
              }
            } else {
              if(handleMessage(&comms[i].incomingMessage, &comms[i].outgoingMessage)){
                writeMessageToAllBut(999, &comms[i].outgoingMessage);
              }
            } 
          }
        }
        while (comms[i].available()){
          comms[i].read();
        }
      }
   }
}
byte coms0Read() {
  return Coms0.read();
}

int coms0Available() {
  return Coms0.available();
}

void coms0Write(char *c, int i) {
  Coms0.write(c, i);
  Coms0.write("\n");
}

byte coms1Read() {
  return Coms1.read();
}

int coms1Available() {
  return Coms1.available();
}

void coms1Write(char *c, int i) {
  Coms1.write(c, i);
  Coms1.write("\n");
}

byte coms3Read() {
  return Coms3.read();
}

int coms3Available() {
  return Coms3.available();
}

void coms3Write(char *c, int i) {
  Coms3.write(c, i);
  Coms3.write("\n");
}

byte coms4Read() {
  return Coms4.read();
}

int coms4Available() {
  return Coms4.available();
}

void coms4Write(char *c, int i) {
  Coms4.write(c, i);
  Coms4.write("\n");
}

byte coms5Read() {
  return Coms5.read();
}

int coms5Available() {
  return Coms5.available();
}

void coms5Write(char *c, int i) {
  Coms5.write(c, i);
  Coms5.write("\n");
}




char voidRead(){return '-';}
bool voidAvailable() {return false;}
void voidWrite(char *c, int i) {};
