#include <ACANFD_DataBitRateFactor.h>
#include <ACAN_T4.h>
#include <ACAN_T4FD_Settings.h>
#include <ACAN_T4_Settings.h>
#include <ACAN_T4_T4FD_rootCANClock.h>
#include <CANFDMessage.h>
#include <CANMessage.h>

IntervalTimer myTimer;
int ledState = LOW;
const int ledPin = LED_BUILTIN;  // the pin with a LED
void setup() {
  // put your setup code here, to run once:
  Serial.begin(115200);
    pinMode(9, INPUT_PULLUP);

  pinMode(10, INPUT_PULLUP);
  pinMode(11, INPUT_PULLUP);
  pinMode(12, INPUT_PULLUP);
  setupTurnout(0,14,15);
  setupTurnout(1,16,17);
  setupTurnout(2,18,19);
  setupTurnout(3,20,21);
  Serial.println("start");
  pinMode(13, OUTPUT);
    myTimer.begin(blinkLED, 150000);  // blinkLED to run every 0.15 seconds
}


// functions called by IntervalTimer should be short, run as quickly as
// possible, and should avoid calling other functions if possible.
void blinkLED() {
  if (ledState == LOW) {
    ledState = HIGH;
     myTimer.begin(blinkLED,150000); 
  } else {
    ledState = LOW;
    myTimer.begin(blinkLED,150000); 
  }
  digitalWrite(ledPin, ledState);
}
int switchPins[4] = {9,10,11,12};
bool switchButtonStates[] = {false, false, false, false} ;
#define CMDSIZE 4

char cmdID;
int cmdSize;
bool readCMD;
void loop() {
  if (Serial.available()){
    cmdID = Serial.read();
    if (cmdID == 48){
      setDirection(0, true);
    } 
    if (cmdID == 49){
      setDirection(0, false);
    }
    if (cmdID == 50){
      setDirection(1, true);
    } 
    if (cmdID == 51){
      setDirection(1, false);
    }

    if (cmdID == 52){
      setDirection(2, true);
    } 
    if (cmdID == 53){
      setDirection(2, false);
    }
    if (cmdID == 54){
      setDirection(3, true);
    } 
    if (cmdID == 55){
      setDirection(3, false);
    }
  }
 
  delay(50);
}

