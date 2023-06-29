#include "message.h"
# define DEBUG_SENDBACK true

typedef byte (*myRead)();
typedef void (*myWrite)(char c);
typedef int (*myAvailable)();


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
