package hwconfig

import (
	"errors"
	"github.com/nanderv/traincontrol-prototype/internal/bridge/domain"
	"log/slog"
	"sync"
)

type node struct {
	Mac  [3]byte
	Addr byte
}
type HwConfigurator struct {
	sync.Mutex
	nodes map[[3]byte]node

	bridges []BridgeSender[domain.Msg]
}

func NewHWConfigurator() *HwConfigurator {
	return &HwConfigurator{nodes: make(map[[3]byte]node)}
}
func (c *HwConfigurator) AddCommandBridge(cc BridgeSender[domain.Msg]) {
	c.bridges = append(c.bridges, cc)
	return
}
func (c *HwConfigurator) firstFreeAddr() byte {
	for i := byte(1); i < 255; i++ {
		v := false
		for _, no := range c.nodes {
			if no.Addr == i {
				v = true
			}
		}
		if !v {
			return i
		}
	}
	panic("Too many nodes")
}

func (c *HwConfigurator) sendToBridges(msg domain.Msg) error {
	var errOut error
	for _, b := range c.bridges {
		err := b.Send(msg)
		if err != nil {
			if errOut == nil {
				errOut = err
			} else {
				errOut = errors.Join(errOut, err)
			}
		}
	}
	return errOut
}
func (c *HwConfigurator) sendNodeNfo(nde node) {
	slog.Info("New Addr", "Mac", nde.Mac, "Addr", nde.Addr)

}

func (c *HwConfigurator) HandleNodeAnnounce(mac [3]byte, prefAddr byte) {
	c.Lock()
	defer c.Unlock()
	foundNode, nde := c.getNodeByMac(mac)

	if !foundNode {
		nde = c.newNodeWithPreferredAddr(mac, prefAddr)
		c.nodes[nde.Mac] = nde
		slog.Info("New Node", "node", nde)
	} else {
		slog.Info("Welcome back", "node", nde)
	}

	if prefAddr != nde.Addr {
		c.sendNodeNfo(nde)
	} else {
		// update to next phase
	}
}

func (c *HwConfigurator) getNodeByMac(mac [3]byte) (bool, node) {
	nde, ok := c.nodes[mac]
	return ok, nde
}

func (c *HwConfigurator) newNodeWithPreferredAddr(mac [3]byte, prefAddr byte) node {
	if prefAddr == 255 {
		return node{Mac: mac, Addr: c.firstFreeAddr()}
	}

	for _, no := range c.nodes {
		if no.Addr == prefAddr && no.Mac != mac {
			return node{Mac: mac, Addr: c.firstFreeAddr()}
		}
	}

	return node{Mac: mac, Addr: prefAddr}
}
