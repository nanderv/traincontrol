#include "message.h"
# define DEBUG_SENDBACK false

typedef byte (*myRead)();
typedef void (*myWrite)(char c);
typedef int (*myAvailable)();
typedef bool (*myHandler)(messageSlot *handleMessage, messageSlot *sendBack);


struct twoHex {
  char fst;
  char snd;
};

typedef twoHex th;

struct comsChannel {
  bool inUse;
  byte bitmask;
  byte halfRead;
  byte isHalfRead;
  messageSlot msg;
  myWrite write;
  myAvailable available;
  myRead read;
};

comsChannel comss[CHAN_IN_USE];
messageSlot outboundMSG[CHAN_IN_USE];

struct handler {
  bool inUse;
  byte typeHandled;
  myHandler handler;
};

handler handlers[16];
bool handleMessage(messageSlot *in, messageSlot *out){
  out -> type = in -> type;
  for(int i=0; i< MSG_LENGTH;i++){
    out->content[i] =in -> content[i];
  }
  
  out ->checkByte = in ->checkByte;
  return true;
}