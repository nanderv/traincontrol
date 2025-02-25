#include "Arduino.h"
#include "hexConv.h"
#include "coms.h"
#include "defines.h"
#include "startup.h"
#include "globals.h"


void writeToAllBut(int dest, char *input) {
    for (int i = 0; i < CHAN_IN_USE; i++) {
        if (comms[i].inUse && (DEBUG_SENDBACK || i != dest)) {
            comms[i].write(input, 2 * (MSG_LENGTH + 2));
        }
    }
}

char messageBuf[2 * (MSG_LENGTH + 2)];

void convertMessageToByte(messageSlot *msg, char *buffer) {
    twoHex inB;
    inB = fromHex(msg->type);
    buffer[0] = inB.fst;
    buffer[1] = inB.snd;

    for (int i = 0; i < MSG_LENGTH; i++) {
        inB = fromHex(msg->content[i]);
        buffer[2 * (i + 1)] = inB.fst;
        buffer[2 * (i + 1) + 1] = inB.snd;
    }
    inB = fromHex(msg->checkByte);
    buffer[2 * (MSG_LENGTH + 1)] = inB.fst;
    buffer[2 * (MSG_LENGTH + 1) + 1] = inB.snd;
}

void writeMessageToAllBut(int dest, messageSlot *msg) {
    convertMessageToByte(msg, messageBuf);
    writeToAllBut(dest, messageBuf);
}


// True if finished, otherwise false.
bool readMessage(int comNR) {
    messageSlot *slot = &(comms[comNR].incomingMessage);
    byte b1 = comms[comNR].read();
    if (toHex(b1) == 128) {
        return false;
    }
    byte b2 = comms[comNR].read();
    if (toHex(b2) == 128) {
        return false;
    }
    if (!toHexDuo(b1, b2, &slot->type)) {
        return false;
    }
    byte checkByte = slot->type;

    for (int i = 0; i < MSG_LENGTH; i++) {
        b1 = comms[comNR].read();
        if (toHex(b1) == 128) {
            return false;
        }
        b2 = comms[comNR].read();
        if (toHex(b2) == 128) {
            return false;
        }
        if (!toHexDuo(b1, b2, &slot->content[i])) {
            return false;
        };
        checkByte = checkByte ^ slot->content[i];
    }
    b1 = comms[comNR].read();
    if (toHex(b1) == 128) {
        return false;
    }
    b2 = comms[comNR].read();
    if (toHex(b2) == 128) {
        return false;
    }
    if (!toHexDuo(b1, b2, &slot->checkByte)) {
        return false;
    }

    return checkByte == slot->checkByte;
}

void handleChannel(int i) {
    if (comms[i].inUse) {
        if (comms[i].available() >= (2 * (MSG_LENGTH + 2))) {
            while (comms[i].available() >= (2 * (MSG_LENGTH + 2))) {
                if (readMessage(i)) {
                    if (comms[i].incomingMessage.type == 0) {
                        if (handleZeroMode(&comms[i].incomingMessage, &comms[i].outgoingMessage)) {
                            writeMessageToAllBut(999, &comms[i].outgoingMessage);
                        }
                    } else {
                        writeMessageToAllBut(i, &(comms[i].incomingMessage));
                        if (handleMessage(&comms[i].incomingMessage, &comms[i].outgoingMessage)) {
                            writeMessageToAllBut(999, &comms[i].outgoingMessage);
                        }
                    }
                }
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
    Coms0.print("\n");
}

byte coms1Read() {
    return Coms1.read();
}

int coms1Available() {
    return Coms1.available();
}

void coms1Write(char *c, int i) {
    Coms1.write(c, i);
    Coms1.print("\n");
}

byte coms2Read() {
    return Coms2.read();
}

int coms2Available() {
    return Coms2.available();
}

void coms2Write(char *c, int i) {
    Coms2.write(c, i);
    Coms2.print("\n");
}


char voidRead() { return '-'; }

bool voidAvailable() { return false; }

void voidWrite(char *c, int i) {};
