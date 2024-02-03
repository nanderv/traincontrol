# define MSG_LENGTH 6
# define NUM_HANDLERS 16



struct messageSlot {
  byte type;
  byte content[MSG_LENGTH];
  byte checkByte;
};

typedef bool (*myHandler)(messageSlot *handleMessage, messageSlot *sendBack);

void setCheckByte(messageSlot *m){
  byte check = m->type;
    for(int i=0; i< MSG_LENGTH;i++){
      check = check ^ m->content[i];
    }
  m->checkByte = check;
}

struct handler {
  bool inUse;
  byte typeHandled;
  myHandler handler;
};

handler handlers[NUM_HANDLERS];
bool handleMessage(messageSlot *in, messageSlot *out){
  for (int i=0; i<NUM_HANDLERS;i++){
    if (handlers[i].inUse){
      if (handlers[i].typeHandled == in->type){
        if (handlers[i].handler(in, out)){
          return true;
        }
      }
    }
  }
  return false;
}

bool addHandler(byte typeHandled, myHandler func){
  for (int i=0;i<NUM_HANDLERS;i++){
    if (!(handlers[i].inUse)){
        handlers[i].inUse = true;
        handlers[i].typeHandled = typeHandled;
        handlers[i].handler=func;
        return true;
    }
  }
  return false;
}
bool echoEcho(messageSlot *handleMessage, messageSlot *sendBack){
    sendBack -> type = handleMessage -> type;
  for(int i=0; i< MSG_LENGTH;i++){
    sendBack->content[i] =handleMessage -> content[i];
  }
  
  sendBack ->checkByte = handleMessage ->checkByte;
  return true;
}

bool flashingLights(messageSlot *handleMessage, messageSlot *sendBack){
  digitalWrite(LED_BUILTIN, HIGH);
  delay(handleMessage->content[0]*10);
  digitalWrite(LED_BUILTIN, LOW);
  sendBack ->type= 2;
    for(int i=0; i< MSG_LENGTH;i++){
    sendBack->content[i] =handleMessage -> content[i];
  }
  setCheckByte(sendBack);
    return false;
}

bool sendAck(messageSlot *handleMessage, messageSlot *sendBack, byte id){
  sendBack->type = id;
  for (int i=0;i<MSG_LENGTH;i++){
    sendBack->content[i] = handleMessage->content[i];
  }
  setCheckByte(sendBack);
  return true;
}