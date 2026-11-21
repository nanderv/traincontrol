#include "Arduino.h"
#include "defines.h"
#include "coms.h"
#include "startup.h"
#include "globals.h"
#include <EEPROM.h>

#include <stdio.h>
#include "pico/stdlib.h"
#include "pico/unique_id.h"


char DEVICE_ID = 0;
void setup() {

    EEPROM.begin(256);
    Coms0.begin(115200);
    Coms1.begin(9600);
    Coms2.begin(9600);

    //   LoadMemory();

  pinMode(LED_BUILTIN, OUTPUT);

  pinMode(2, OUTPUT);
  pinMode(3, OUTPUT);
  pinMode(4, OUTPUT);
  pinMode(5, OUTPUT);

  pinMode(9, OUTPUT);
  pinMode(10, OUTPUT);
  pinMode(11, OUTPUT);
  pinMode(12, OUTPUT);
  comms[0].write = coms0Write;
  comms[0].available = coms0Available;
  comms[0].read = coms0Read;
  comms[0].inUse = true;
  comms[1].write = coms1Write;
  comms[1].available = coms1Available;
  comms[1].read = coms1Read;
  comms[1].inUse = true;
  comms[2].write = coms2Write;
  comms[2].available = coms2Available;
  comms[2].read = coms2Read;
  comms[2].inUse = true;
}



struct LocDef {
  char addr[2 * PICO_UNIQUE_BOARD_ID_SIZE_BYTES + 1];
  char my_id;
};


#define ADDRESSES_LEN = 1
constexpr LocDef Addresses[] = {
  {
    "E66164084342B537",
    1
  },
{
  "E66164084342B599",
  2
},
};


void loop() {

  if (DEVICE_ID > 0) {
    static int sleepCount;
    for (int i = 0; i < CHAN_IN_USE; i++) {
      handleChannel(i);
    }
    if(sleepCount ==0){
      digitalWrite(LED_BUILTIN, HIGH);
      handleStartState();
    }

    if(sleepCount == 1000){
      digitalWrite(LED_BUILTIN, LOW);
    }

    sleepCount++;
    if (sleepCount == 5000){
      sleepCount = 0;
    }

    delay(1);
    Serial.flush();
  } else {

    char id_str[2 * PICO_UNIQUE_BOARD_ID_SIZE_BYTES + 1];
    pico_get_unique_board_id_string(id_str, sizeof(id_str));
    for (int i = 0; i < 1; i++) {
      if (strcmp( Addresses[i].addr, id_str)==0){
        DEVICE_ID = Addresses[i].my_id;
        return;
      }
    }

    Serial.printf("Unique ID: %s\n", id_str);
    delay(2000);
  }
}
