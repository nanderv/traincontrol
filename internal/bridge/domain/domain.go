package domain

import (
	"errors"
)

const hextable = "0123456789abcdef"
const reverseHexTable = "" +
	"\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff" +
	"\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff" +
	"\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff" +
	"\x00\x01\x02\x03\x04\x05\x06\x07\x08\x09\xff\xff\xff\xff\xff\xff" +
	"\xff\x0a\x0b\x0c\x0d\x0e\x0f\xff\xff\xff\xff\xff\xff\xff\xff\xff" +
	"\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff" +
	"\xff\x0a\x0b\x0c\x0d\x0e\x0f\xff\xff\xff\xff\xff\xff\xff\xff\xff" +
	"\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff" +
	"\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff" +
	"\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff" +
	"\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff" +
	"\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff" +
	"\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff" +
	"\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff" +
	"\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff" +
	"\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff"

const Size = 6

func ValidChar(b byte) bool {
	rr := reverseHexTable[b]
	return rr <= 0x0f

}

type Msg struct {
	Type byte
	Val  [Size]byte
}

func (m Msg) Encode() (r RawMsg) {
	cd := m.Type
	r[0] = hextable[m.Type>>4]
	r[1] = hextable[m.Type&0x0f]
	j := 2
	for _, v := range m.Val {
		cd ^= v
		r[j] = hextable[v>>4]
		r[j+1] = hextable[v&0x0f]
		j += 2
	}
	r[j] = hextable[cd>>4]
	r[j+1] = hextable[cd&0x0f]

	return r
}

func decode(dst *[Size + 2]byte, src [(Size + 2) * 2]byte) (int, error) {
	i, j := 0, 1
	for ; j < len(src); j += 2 {
		p := src[j-1]
		q := src[j]

		a := reverseHexTable[p]
		b := reverseHexTable[q]
		if a > 0x0f {
			return j - 1, errors.New("oops")
		}
		if b > 0x0f {
			return j, errors.New("oops")
		}
		dst[i] = (a << 4) | b
		i++
	}
	return j, nil
}

const RawSize = (Size + 2) * 2

type RawMsg [RawSize]byte

func (m RawMsg) String() string {
	return string(m[:])
}

func (m RawMsg) Decode() (Msg, error) {
	var r Msg
	var bytes [Size + 2]byte
	_, err := decode(&bytes, m)
	var chk byte
	if err != nil {
		return Msg{}, err
	}
	for i, b := range bytes {
		chk ^= b
		if i == 0 {

			r.Type = b
		} else {
			if i < Size+1 {
				r.Val[i-1] = b
			}
		}
	}
	if chk != 0 {
		return Msg{}, errors.New("incorrect check number")
	}
	return r, nil
}