package hwconfig

import "github.com/nanderv/traincontrol-prototype/internal/hwconfig/domain/node"

type BridgeSender interface {
	SendNodeInfoUpdate(node node.Node, state byte) error
	SendEepromWrite(nodeAddr byte, commandSlot byte, typ byte, data [2]byte) error
}
