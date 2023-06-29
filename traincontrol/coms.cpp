#include "coms.h"
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

void writeMessageToAllBut(int dest, MessageSlot *msg) {
  char checkByte;
  checkByte = 0;
  writeToAllBut(dest, msg->type);
  checkByte = checkByte ^ msg->type;
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
bool readByte(char in, int comNR) {
  MessageSlot *slot = &(comss[comNR].msg);

  (*slot).checkByte = slot->checkByte ^ in;
  switch ((*slot).readState) {
    case type:
      (*slot).type = in;
      (*slot).checkByte = in;
      (*slot).readState = content;
      (*slot).readCounter = 0;
      break;
    case content:

      (*slot).content[(*slot).readCounter] = in;
      (*slot).readCounter++;
      if ((*slot).readCounter >= MSG_LENGTH) {
        (*slot).readCounter = 0;
        (*slot).readState = check;
      }
      break;
    case check:
      (*slot).readState = type;
      if ((*slot).checkByte == 0) {
        writeMessageToAllBut(comNR, slot);
        (*slot).checkByte = 0;
        Serial.println("!");

        return true;
      } else {
        Serial.println("INVALID");
      }
  }
  return false;
}

void handleChannel(int i) {
  if (comss[i].inUse) {
    while (comss[i].available()) {

      char in = comss[i].read();

      char hexV = toHex(in);
      if (hexV < 127) {
        // Finish read
        if (comss[i].isHalfRead) {
          char inByte = (comss[i].halfRead) * 16 + hexV;
          if (readByte(inByte, i)) {
            break;
          }
          comss[i].isHalfRead = false;
          comss[i].halfRead = 0;
        } else {
          // Store that we read half a byte
          comss[i].isHalfRead = true;
          comss[i].halfRead = hexV;
        }
      } else {
        // We should reset if we read an incorrect byte.
        resetReader(i);
        Coms1.println("Reset");
        break;
      }
    }
  }
}


char coms1Read() {
  return Coms1.read();
}

bool coms1Available() {
  return Coms1.available();
}

void coms1Write(char c) {
  Coms1.write(c);
}

char voidRead() { return '-'; }

bool voidAvailable() { return false; }

void voidWrite(char c) {};
