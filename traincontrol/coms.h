#include "message.h"

# define DEBUG_SENDBACK true
# define CHAN_IN_USE 1
#define Coms1 Serial

typedef char (*myRead)();

typedef void (*myWrite)(char c);

typedef bool (*myAvailable)();


struct twoHex {
    char fst;
    char snd;
};

typedef twoHex th;

struct comsChannel {
    bool inUse;
    char bitmask;
    char halfRead;
    char isHalfRead;
    MessageSlot msg;
    Serial serial;
    myWrite write;
    myAvailable available;
    myRead read;
};

comsChannel comss[CHAN_IN_USE];

void writeToAllBut(int dest, char input);

void writeMessageToAllBut(int dest, MessageSlot *msg);
void resetReader(int dest);
bool readByte(char in, int comNR);
void handleChannel(int i);
char coms1Read();
bool coms1Available();
void coms1Write(char c);
char voidRead();
bool voidAvailable();
void voidWrite(char c);
