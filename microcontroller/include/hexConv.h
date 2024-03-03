#ifndef TR_HEXCONV_H
#define TR_HEXCONV_H

#include "Arduino.h"


struct twoHex {
    char fst;
    char snd;
};

byte toHex(byte hv);

twoHex fromHex(byte hv);

bool toHexDuo(byte hi, byte lo, byte *rr);

#endif //TR_HEXCONV_H
