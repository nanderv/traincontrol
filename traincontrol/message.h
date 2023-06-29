# define MSG_LENGTH 6

enum Mode {
  type,
  content,
  check,
};

struct MessageSlot {
  char type;
  char content[MSG_LENGTH];
  char checkByte;
  Mode readState;
  int readCounter;
};
