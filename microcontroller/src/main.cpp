#include <PN532_I2C.h>
#include <PN532.h>


PN532_I2C pn532(Wire1);
#define PN532IRQPIN D8
PN532_I2C pn532B(Wire);
#define PN532IRQPINB D2
volatile boolean cState = false;
volatile boolean cState2 = false;

PN532 nfc(pn532);
PN532 nfcB(pn532B);

void handleA();

void handleB();

void cardreading();

void cardreading2();

void setup() {
    bool setSDA(D4);
    bool setSCL(D5);
    delay(1000);
    // put your setup code here, to run once:
    Serial.begin(115200);
    Serial.println("\nHello!");
    Wire1.setSDA(D6);
    Wire1.setSCL(D7);
    delay(100);
    pn532 = PN532_I2C(Wire1);
    nfc = PN532(pn532);
    nfc.begin();
    uint32_t versiondata = nfc.getFirmwareVersion();
    if (!versiondata) {
        Serial.print("Didn't find PN53x bread");
        Serial.flush();
        while (1); // halt
    }


    nfcB.begin();
    versiondata = nfcB.getFirmwareVersion();
    if (!versiondata) {
        Serial.print("Didn't find PN53x bread");
        Serial.flush();
        while (1); // halt
    }
    // Got ok data, print it out!
    Serial.print("Found chip PN5");
    Serial.println((versiondata >> 24) & 0xFF, HEX);
    Serial.print("Firmware ver. ");
    Serial.print((versiondata >> 16) & 0xFF, DEC);
    Serial.print('.');
    Serial.println((versiondata >> 8) & 0xFF, DEC);

    Serial.println(cState);
    attachInterrupt(digitalPinToInterrupt(PN532IRQPIN), cardreading, FALLING);
    attachInterrupt(digitalPinToInterrupt(PN532IRQPINB), cardreading2, FALLING);
    //It generates interrupt, I do not really know why?!
    nfc.SAMConfig();
    nfcB.SAMConfig();

}

void loop() {
    delay(10);

    handleA();
    handleB();
}

volatile int ac = 0;
volatile int bc = 0;

void handleA() {
    // put your main code here, to run repeatedly:
    if (ac > 0) {
        ac--;

    }

    if (cState && ac <= 0) {
        Serial.println("Interrupted A");
        uint8_t success;
        uint8_t uid[] = {0, 0, 0, 0, 0, 0, 0};  // Buffer to store the returned UID
        uint8_t uidLength;                        // Length of the UID (4 or 7 bytes depending on ISO14443A card type)
        // Wait for an ISO14443A type cards (Mifare, etc.).  When one is found
        // 'uid' will be populated with the UID, and uidLength will indicate
        // if the uid is 4 bytes (Mifare Classic) or 7 bytes (Mifare Ultralight)
        success = nfc.readPassiveTargetID(PN532_MIFARE_ISO14443A, uid, &uidLength, 200);

        if (success) {
            // Display some basic information about the card
            Serial.println("Found an ISO14443A card");
            Serial.print("  UID Length: ");
            Serial.print(uidLength, DEC);
            Serial.println(" bytes");
            Serial.print("  UID Value: ");
            nfc.PrintHex(uid, uidLength);

        }
        //This must be called or IRQ won't work!
        nfc.startRead(PN532_MIFARE_ISO14443A, uid, &uidLength);
        cState = false;
        Serial.print("OUTA : ");
        Serial.println(cState);
        ac = 100;
    }
}

void handleB() {
    if (bc > 0) {
        bc--;
    }

    // put your main code here, to run repeatedly:
    if (cState2 && bc <= 0) {

        Serial.println("Interrupted B");


        uint8_t success;
        uint8_t uid[] = {0, 0, 0, 0, 0, 0, 0};  // Buffer to store the returned UID
        uint8_t uidLength;                        // Length of the UID (4 or 7 bytes depending on ISO14443A card type)
        // Wait for an ISO14443A type cards (Mifare, etc.).  When one is found
        // 'uid' will be populated with the UID, and uidLength will indicate
        // if the uid is 4 bytes (Mifare Classic) or 7 bytes (Mifare Ultralight)
        success = nfcB.readPassiveTargetID(PN532_MIFARE_ISO14443A, uid, &uidLength, 200);

        if (success) {
            // Display some basic information about the card
            Serial.println("Found an ISO14443A card");
            Serial.print("  UID Length: ");
            Serial.print(uidLength, DEC);
            Serial.println(" bytes");
            Serial.print("  UID Value: ");
            nfcB.PrintHex(uid, uidLength);

        }
        //This must be called or IRQ won't work!
        nfcB.startRead(PN532_MIFARE_ISO14443A, uid, &uidLength);
        cState2 = false;
        Serial.print("OUT B: ");
        Serial.println(cState2);
        bc = 100;
    }
}

void cardreading() {
        cState = true;
}

void cardreading2() {
        cState2 = true;
}