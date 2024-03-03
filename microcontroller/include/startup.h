
#ifndef TR_STARTUP_H
#define TR_STARTUP_H

#include "message.h"

bool handleZeroMode(messageSlot *handleMessage, messageSlot *sendBack);
void handleStartState();
bool restart(messageSlot *handleMessage, messageSlot *sendBack);

#endif //TR_STARTUP_H
