# define MSG_LENGTH 6

enum mode {
  type,
  content,
  check,
};

struct messageSlot {
  byte type;
  byte content[MSG_LENGTH];
  byte checkByte;
  mode readState;
  int readCounter;
};
