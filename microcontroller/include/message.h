
#ifndef TR_MESSAGE_H
#define TR_MESSAGE_H

#include "defines.h"
struct messageSlot {
    byte type;
    byte content[MSG_LENGTH];
    byte checkByte;
};

typedef bool (*myHandler)(messageSlot *handleMessage, messageSlot *sendBack);

void setCheckByte(messageSlot *m);

struct handler {
    bool inUse;
    byte typeHandled;
    myHandler handler;
};


bool handleMessage(messageSlot *in, messageSlot *out);

bool addHandler(byte typeHandled, myHandler func);

bool echoEcho(messageSlot *handleMessage, messageSlot *sendBack);

bool flashingLights(messageSlot *handleMessage, messageSlot *sendBack);

bool sendAck(messageSlot *handleMessage, messageSlot *sendBack, byte id);
#endif
