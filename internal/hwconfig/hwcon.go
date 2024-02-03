package hwconfig

import (
	"github.com/nanderv/traincontrol-prototype/internal/hwconfig/domain/node"
	"log/slog"
	"sync"
)

type HwConfigurator struct {
	sync.Mutex
	nodes  map[[3]byte]node.Node
	bridge BridgeSender
}

func NewHWConfigurator() *HwConfigurator {
	return &HwConfigurator{nodes: make(map[[3]byte]node.Node)}
}
func (c *HwConfigurator) SetBridgeAdapter(cc BridgeSender) {
	c.bridge = cc
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

func (c *HwConfigurator) sendNodeNfo(nde node.Node) {
	slog.Info("New Addr", "Mac", nde.Mac, "Addr", nde.Addr)
	err := c.bridge.SendNodeInfoUpdate(nde)
	if err != nil {
		slog.Error("Could not send node info", "node", nde, "err", err)
		return
	}
}

func (c *HwConfigurator) HandleNodeAnnounce(mac [3]byte, prefAddr byte) {
	c.Lock()
	defer c.Unlock()
	foundNode, nde := c.getNodeByMac(mac)

	if !foundNode {
		nde = c.newNodeWithPreferredAddr(mac, prefAddr)
		c.nodes[nde.Mac] = nde
		slog.Info("New Node", "Node", nde, "requested", prefAddr)
	} else {
		slog.Info("Welcome back", "Node", nde, "requested", prefAddr)
	}

	if prefAddr != nde.Addr || !foundNode {
		c.sendNodeNfo(nde)
	} else {

		// update to next phase
	}
}

func (c *HwConfigurator) getNodeByMac(mac [3]byte) (bool, node.Node) {
	nde, ok := c.nodes[mac]
	return ok, nde
}

func (c *HwConfigurator) newNodeWithPreferredAddr(mac [3]byte, prefAddr byte) node.Node {
	if prefAddr == 255 {
		return node.Node{Mac: mac, Addr: c.firstFreeAddr()}
	}

	for _, no := range c.nodes {
		if no.Addr == prefAddr && no.Mac != mac {
			return node.Node{Mac: mac, Addr: c.firstFreeAddr()}
		}
	}

	return node.Node{Mac: mac, Addr: prefAddr}
}
