void setup() {
  // put your setup code here, to run once:
    Serial.begin(115200);
    pinMode(12, INPUT_PULLUP);
}

void loop() {
  // put your main code here, to run repeatedly:
  if (digitalRead(12)) {
    // do this if C2 is high
    Serial.println("Hi");
  } else {
    // do this if C2 is low
    Serial.println("Lo");
  }
  delay(1000);
}
