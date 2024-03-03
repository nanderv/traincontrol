#include "Arduino.h"
#include "hexConv.h"


byte toHex(byte hv){
  switch (hv) {
    case '0': return 0;
    case '1': return 1;
    case '2': return 2;
    case '3': return 3;
    case '4': return 4;
    case '5': return 5;
    case '6': return 6;
    case '7': return 7;
    case '8': return 8;
    case '9': return 9;
    
    case 'a':
    case 'A': return 10;
    
    case 'b':
    case 'B': return 11;

    case 'c':
    case 'C': return 12;
    
    case 'd':
    case 'D': return 13;

    case 'e':
    case 'E': return 14;

    case 'f':
    case 'F': return 15;
  }
  return 128;
}

// A4 -> 10*16 + 4
// R# -> ERROR
bool toHexDuo(byte hi, byte lo, byte *rr ){
  byte result = toHex(hi);
  if (result > 16){
    return false;
  }
  byte loResult = toHex(lo);
  if (result > 16){
    return false ;
  }
  *rr =  result*16+loResult;
  return true;
}

char fromHexSub(byte hv){
  switch(hv){
    case 0: return '0';
    case 1: return '1';
    case 2: return '2';
    case 3: return '3';
    case 4: return '4';
    case 5: return '5';
    case 6: return '6';
    case 7: return '7';
    case 8: return '8';
    case 9: return '9';
    case 10: return 'a';
    case 11: return 'b';
    case 12: return 'c';
    case 13: return 'd';
    case 14: return 'e';
    case 15: return 'f';
    default: return '-';
  }
}

twoHex fromHex(byte hv){
  twoHex result;
  result.fst = fromHexSub(hv/16);
  result.snd = fromHexSub(hv%16);
  return result;
}
