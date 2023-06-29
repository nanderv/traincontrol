# define MSG_LENGTH 6
# define NUM_HANDLERS 16
enum mode {
  type,
  content,
  check,
};

struct messageSlot {
  byte type;
  byte content[MSG_LENGTH];
  byte checkByte;
  mode readState;
  int readCounter;
};
typedef bool (*myHandler)(messageSlot *handleMessage, messageSlot *sendBack);


messageSlot outboundMSG[CHAN_IN_USE];

struct handler {
  bool inUse;
  byte typeHandled;
  myHandler handler;
};

handler handlers[NUM_HANDLERS];
bool handleMessage(messageSlot *in, messageSlot *out){
  out -> type = in -> type;
  for(int i=0; i< MSG_LENGTH;i++){
    out->content[i] =in -> content[i];
  }
  
  out ->checkByte = in ->checkByte;
  return true;
}

bool addHandler(byte typeHandled, myHandler func){
  for (int i=0;i<NUM_HANDLERS;i++){
    if (!(handlers[i].inUse)){
        handlers[i].inUse = true;
        handlers[i].typeHandled = type;
        handlers[i].handler=func;
        return true;
    }
  }
  return false;
}