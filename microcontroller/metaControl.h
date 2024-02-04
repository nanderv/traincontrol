#define MAX_SLOTS 32
#include <EEPROM.h>

struct ControlSlot{
    byte id; // THE ID in the array;
    byte type;
    byte data[2];
};

byte MY_ID;
ControlSlot controls[MAX_SLOTS];

bool metaControlWriteEEPROM(messageSlot *handleMessage, messageSlot *sendBack){
    if (handleMessage->content[1] != MY_ID){
        return false;
    }
    
    if (handleMessage->content[2] > (MAX_SLOTS + 1) ){
        sendBack->type=255;
        sendBack->content[0]=handleMessage->type;  
        setCheckByte(sendBack);
        return true;
    }

    int slot = handleMessage->content[2]+1;
    if (EEPROM.read(slot) != handleMessage->content[3]){
        EEPROM.write(slot, handleMessage->content[3]);
    }
    if (EEPROM.read(slot+1) != handleMessage->content[4]){
        EEPROM.write(slot+1, handleMessage->content[4]);
    }
    if (EEPROM.read(slot+2) != handleMessage->content[5]){
        EEPROM.write(slot+2, handleMessage->content[5]);
    }
    sendBack->type=0;
    sendBack->content[0]=EEPROM_WRITE_RETURN;
    sendBack->content[1]=MY_ID;
    sendBack->content[2] = handleMessage->content[2];
    sendBack->content[3] = handleMessage->content[3];
    sendBack->content[4] = handleMessage->content[4];
    sendBack->content[5] = handleMessage->content[5];
    if (handleMessage->content[2] > 0 ){
        int slot = (handleMessage->content[2]+1)*4;
        byte i = handleMessage->content[2];
        byte type = EEPROM.read(slot);

        if(type == 255 ){
            return false;
        }
        if(type == 1){
            pinMode(controls[i].data[1], OUTPUT);
            pinMode(controls[i].data[1]+1, OUTPUT);
            digitalWrite(controls[i].data[1], LOW);
            digitalWrite(controls[i].data[1]+1, LOW);
        }
        if(type ==2){
            pinMode(controls[i].data[1], OUTPUT);
            digitalWrite(controls[i].data[1], HIGH);
        }
    }

    setCheckByte(sendBack);

    return true;
}

bool metaControlReadEEPROM(messageSlot *handleMessage, messageSlot *sendBack){
    int slot = handleMessage->content[1]*4;
    sendBack->type=0;
    sendBack->content[0]=EEPROM_READ_RETURN;
    sendBack->content[1]=MY_ID;

    sendBack->content[2] = EEPROM.read(slot);
    sendBack->content[3] = EEPROM.read(slot+1);
    sendBack->content[4] = EEPROM.read(slot+2);
    setCheckByte(sendBack);
    return true;
}

bool LoadMemory(){
    MY_ID = EEPROM.read(0);
    for ( int i = 0; i < MAX_SLOTS; i++){
        int slot = (i+1)*4;
        byte type = EEPROM.read(slot);
        controls[i].type = type;
        controls[i].data[0]= EEPROM.read(slot+1);
        controls[i].data[1]= EEPROM.read(slot+2);

        if(type == 255 ){
            continue;
        }
        if(type == 1){
            pinMode(controls[i].data[1], OUTPUT);
            pinMode(controls[i].data[1]+1, OUTPUT);
            digitalWrite(controls[i].data[1], LOW);
            digitalWrite(controls[i].data[1]+1, LOW);
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
