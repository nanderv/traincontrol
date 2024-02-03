#include "message.h"
# define DEBUG_SENDBACK false

// Wrappers for Serial.read, write, available.
typedef byte (*myRead)();
typedef void (*myWrite)(char *c, int length);
typedef int (*myAvailable)();


struct twoHex {
  char fst;
  char snd;
};

typedef twoHex th;

struct comsChannel {
  bool inUse; // Is this channel used? If false, ignore it entirely
  messageSlot incomingMessage; // 
  messageSlot outgoingMessage; // The message to send back

  // Functions for serial functionality
  myWrite write;
  myAvailable available;
  myRead read;
};

comsChannel comms[CHAN_IN_USE];
