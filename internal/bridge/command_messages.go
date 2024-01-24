package bridge

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

type Msg struct {
	Type byte
	Val  [Size]byte
}

func (m Msg) encode() (r rawMsg) {
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
			return i, errors.New("oops")
		}
		if b > 0x0f {
			return i, errors.New("oops")
		}
		dst[i] = (a << 4) | b
		i++
	}
	if len(src)%2 == 1 {
		// Check for invalid char before reporting bad length,
		// since the invalid char (if present) is an earlier problem.
		if reverseHexTable[src[j-1]] > 0x0f {
			return i, errors.New("oops")
		}
		return i, nil
	}
	return i, nil
}

type rawMsg [(Size + 2) * 2]byte

func (m rawMsg) String() string {
	return string(m[:])
}

func (m rawMsg) decode() (r Msg, err error) {
	var bytes [Size + 2]byte
	_, err = decode(&bytes, m)
	var chk byte
	if err != nil {
		return
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
		err = errors.New("incorrect check number")
	}
	return
}
