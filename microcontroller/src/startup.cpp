#include <EEPROM.h>
#include "globals.h"
#include "coms.h"
#include "startup.h"
#include <EEPROM.h>


void handleStartState() {
    static messageSlot startupSendSlot;

    startupSendSlot.type = 0;
    startupSendSlot.content[0] = 2;
    startupSendSlot.content[1] = DEVICE_ID;

    setCheckByte(&startupSendSlot);

    writeMessageToAllBut(999, &startupSendSlot);
}


bool restart(messageSlot *handleMessage, messageSlot *sendBack) {
    if (handleMessage->content[1] == MAC[0] && handleMessage->content[2] == MAC[1] &&
        handleMessage->content[3] == MAC[2]) {
        delay(2000);
    }
    return false;
}
