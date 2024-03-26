package datasets

import (
	"github.com/nanderv/traincontrol-prototype/datasets/test"
	hwDomain "github.com/nanderv/traincontrol-prototype/internal/hardware/domain"
)

func GetHardWareStateByID(stateID string) hwDomain.HardwareState {
	if stateID == "test" {
		return test.GetBaseLayout()
	}
	panic("Incorrect state")
	return hwDomain.NewHardwareState()
}
