#include <EEPROM.h>
int ADDRESS_EEPROM = 0;
byte MY_ADDR;

void write_device_address(byte address) {
   EEPROM.write(ADDRESS_EEPROM, address);
}

byte  get_device_address(byte address){
  MY_ADDR  = EEPROM.read(ADDRESS_EEPROM);
  return MY_ADDR;
}