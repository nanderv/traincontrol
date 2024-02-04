#define MAX_SLOTS 32
#include <EEPROM.h>

bool metaControlWriteEEPROM(messageSlot *handleMessage, messageSlot *sendBack){
  int slot = handleMessage->content[1]*4;

  if (handleMessage->content[1] > (MAX_SLOTS + 1) ){
    sendBack->type=255;
    sendBack->content[0]=handleMessage->type;  
    setCheckByte(sendBack);
    return true;
  }

  EEPROM.write(slot, handleMessage->content[2]);
  EEPROM.write(slot+1, handleMessage->content[3]);
  EEPROM.write(slot+2, handleMessage->content[4]);
  EEPROM.write(slot+3, handleMessage->content[5]);
  return false;
}

bool metaControlReadEEPROM(messageSlot *handleMessage, messageSlot *sendBack){
    int slot = handleMessage->content[1]*4;
    sendBack->type=255;
    sendBack->content[0]=handleMessage->type;
    sendBack->content[1]=handleMessage->content[1];

    sendBack->content[2] = EEPROM.read(slot);
    sendBack->content[3] = EEPROM.read(slot+1);
    sendBack->content[4] = EEPROM.read(slot+2);
    sendBack->content[5] = EEPROM.read(slot+3);
    setCheckByte(sendBack);
    return true;
}
bool metaControl(messageSlot *handleMessage, messageSlot *sendBack){
    if (handleMessage->content[0] == 254){
        return metaControlWriteEEPROM(handleMessage, sendBack);
    }
    if(handleMessage->content[0] == 255){
        return metaControlReadEEPROM(handleMessage, sendBack);
    }
    return false;
}


byte MY_ID;

struct ControlSlot{
    byte id; // THE ID in the array;
    byte type;
    byte data[3];
};


ControlSlot controls[MAX_SLOTS];

bool LoadMemory(){
    MY_ID = EEPROM.read(0);
    for ( int i = 0; i < MAX_SLOTS; i++){
        int slot = (i+1)*4;
        byte type = EEPROM.read(slot);

        if(type != 255 ){
            controls[i].id=i;
            controls[i].type=type;
            controls[i].data[0] = EEPROM.read(slot+1);
            controls[i].data[1]  = EEPROM.read(slot+2);
            controls[i].data[2]  = EEPROM.read(slot+3);
        }
        if(type == 1){
            pinMode(controls[i].data[1], OUTPUT);
            pinMode(controls[i].data[2], OUTPUT);
        }
        if(type ==2){
            pinMode(controls[i].data[1], OUTPUT);
            digitalWrite(controls[i].data[1], HIGH);
        }
        delay(1);
    }
    return false;
}



bool flashingLights(messageSlot *handleMessage, messageSlot *sendBack){
    if (handleMessage->content[0] != MY_ID && handleMessage -> content[0] != 255){
        return false;
    }
    digitalWrite(LED_BUILTIN, HIGH);
    delay(handleMessage->content[1]*10);
    digitalWrite(LED_BUILTIN, LOW);
    return sendAck(handleMessage, sendBack, 2);
}
