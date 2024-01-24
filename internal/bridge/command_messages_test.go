package bridge

import (
	"github.com/gofiber/fiber/v2/utils"
	"testing"
)

func must[T any](a T, e error) T {
	if e != nil {
		panic(e)
	}
	return a
}
func mustnot[T any](a T, e error) T {
	if e == nil {
		panic("not error")
	}
	return a
}
func TestSwitchMsg(t *testing.T) {
	msg := Msg{Type: 2, Val: [6]byte{}}
	//utils.AssertEqual(t, msg.encode(), true)
	utils.AssertEqual(t, must(msg.encode().decode()), msg)
	v := msg.encode()
	v[0] = v[0] + 1
	mustnot(v.decode())
}
