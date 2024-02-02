package hwconfig

import (
	"github.com/nanderv/traincontrol-prototype/internal/bridge/domain"
)

type node struct {
	mac  [3]byte
	addr byte
}
type HwConfigurator struct {
	nodes []node

	bridges []MessageSender[domain.Msg]
}

func (c *HwConfigurator) AddCommandBridge(cc MessageSender[domain.Msg]) {
	c.bridges = append(c.bridges, cc)
	return
}
