package hwconfig

import "github.com/nanderv/traincontrol-prototype/internal/hwconfig/domain/node"

type BridgeSender interface {
	SendNodeInfoUpdate(node node.Node) error
}
