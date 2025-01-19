#define NFC_INTERFACE_HSU
#include "Arduino.h"
#include <PN532_HSU.h>
#include <PN532.h>
#include <HardwareSerial.h>

PN532_HSU pn532hsu(Serial1);
PN532 nfc(pn532hsu);

void setup(void)
{
    delay(2000);
    nfc.begin();
    Serial.println("HI");

    uint32_t versiondata = nfc.getFirmwareVersion();
    if (! versiondata) {
        Serial.print("Didn't find PN53x board");
        while (1); // halt
    }
    Serial.print("Found chip PN5"); Serial.println((versiondata>>24) & 0xFF);
    Serial.print("Firmware ver. "); Serial.print((versiondata>>16) & 0xFF);
    Serial.print('.'); Serial.println((versiondata>>8) & 0xFF);

    nfc.setPassiveActivationRetries(0xFF);
}

void loop() {
    readNFC();
    digitalWrite(ledpin1, HIGH);
    delay(10);
    digitalWrite(ledpin1, LOW);

}

void readNFC() {
    if (nfc.tagPresent(100))
    {
        NfcTag tag = nfc.read();
        // tag.print();
        tagId = tag.getUidString();
        // Serial.println("Tag id");
        Serial.println(tagId);
    } else {
        Serial.print(".");
    }
}