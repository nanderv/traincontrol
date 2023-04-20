struct directionalSector{
  int pin;
  bool isReversed;
};

struct poweredSector{
  int partOfDirectional;
  int controlPin;
  bool isEnabled;
};

