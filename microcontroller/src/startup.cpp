#include <EEPROM.h>
#include "globals.h"
#include "coms.h"
#include "startup.h"
#include <EEPROM.h>

char addr1 = 0;

void handleStartState() {
    static messageSlot startupSendSlot;
    addr1 = EEPROM.read(0);

    startupSendSlot.type = 0;
    startupSendSlot.content[0] = 2;
    startupSendSlot.content[1] = addr1;

    setCheckByte(&startupSendSlot);

    writeMessageToAllBut(999, &startupSendSlot);
}

bool handleZeroMode(messageSlot *handleMessage, messageSlot *sendBack) {
    if (handleMessage->content[0] == 1) {
        EEPROM.write(0, handleMessage->content[1]);
        addr1=handleMessage->content[1];
        sendBack->type = 0;
        sendBack->content[0] = 1;
        sendBack->content[1] = addr1;

        EEPROM.commit();
        return true;
    }
    if (handleMessage->content[0] == 254) {
        return restart(handleMessage, sendBack);
    }
    return false;
}


bool restart(messageSlot *handleMessage, messageSlot *sendBack) {
    if (handleMessage->content[1] == MAC[0] && handleMessage->content[2] == MAC[1] &&
        handleMessage->content[3] == MAC[2]) {
        delay(2000);
    }
    return false;
}
