package types

import (
	"fmt"
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
	s := SetSwitch{
		SwitchID:  2,
		Direction: false,
	}
	msg := s.ToBridgeMsg()

	//utils.AssertEqual(t, msg.encode(), true)
	fmt.Println(msg.encode())
	utils.AssertEqual(t, must(msg.encode().decode()), msg)
	v := msg.encode()
	v[0] = v[0] + 1
	mustnot(v.decode())
}
