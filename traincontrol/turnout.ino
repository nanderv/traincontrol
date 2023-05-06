struct turnout{
  int hi;
  int lo;
  bool dir;
} turnouts[4];

void setupTurnout(int nr, int hi, int lo){
  turnouts[nr].hi = hi;
  turnouts[nr].lo = lo;
  pinMode(hi, OUTPUT);
  pinMode(lo, OUTPUT);

}
void setDirection(int nr, bool dir){
  if ( turnouts[nr].dir  == dir){
    return;
  }
  if (dir){
  Serial.print(">");
  } else {
    Serial.print("<");
  }
  Serial.println(nr);
  int pin = turnouts[nr].lo;
  if (turnouts[nr].dir){
    pin = turnouts[nr].hi;
  }
  
  digitalWrite(pin, HIGH);
  delay(80);
  digitalWrite(pin, LOW);
  turnouts[nr].dir = dir;
}
void changeSwitchDirection(int nr){
 setDirection(nr, !turnouts[nr].dir);
}