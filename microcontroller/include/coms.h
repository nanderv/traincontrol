
#ifndef TR_COMS_H
#define TR_COMS_H


#include "message.h"
#include "defines.h"

// Wrappers for Serial.read, write, available.
typedef byte (*myRead)();
typedef void (*myWrite)(char *c, int length);
typedef int (*myAvailable)();



struct comsChannel {
  bool inUse; // Is this channel used? If false, ignore it entirely
  messageSlot incomingMessage; //
  messageSlot outgoingMessage; // The message to send back

  // Functions for serial functionality
  myWrite write;
  myAvailable available;
  myRead read;
};


void writeToAllBut(int dest, char *input);
void convertMessageToByte(messageSlot *msg, char *buffer);
void writeMessageToAllBut(int dest, messageSlot *msg);
bool readMessage(int comNR);
void handleChannel(int i);
byte coms0Read();
int coms0Available();
void coms0Write(char *c, int i);
byte coms1Read();
int coms1Available();
void coms1Write(char *c, int i);
byte coms2Read();
int coms2Available();
void coms2Write(char *c, int i);
byte coms4Read();
int coms4Available();
void coms4Write(char *c, int i);
byte coms5Read();
int coms5Available();
void coms5Write(char *c, int i);
char voidRead();
bool voidAvailable();
void voidWrite(char *c, int i);

#endif
