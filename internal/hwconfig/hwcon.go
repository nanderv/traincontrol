package hwconfig

import (
	"errors"
	"fmt"
	"github.com/nanderv/traincontrol-prototype/internal/bridge/domain"
	"log/slog"
)

type node struct {
	mac  [3]byte
	addr byte
}
type HwConfigurator struct {
	nodes []node

	bridges []BridgeSender[domain.Msg]
}

func (c *HwConfigurator) AddCommandBridge(cc BridgeSender[domain.Msg]) {
	c.bridges = append(c.bridges, cc)
	return
}
func (c *HwConfigurator) firstFreeAddr() byte {
	for i := byte(1); i < 255; i++ {
		v := false
		for _, no := range c.nodes {
			if no.addr == i {
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
func (c *HwConfigurator) sendNodeNfo(mac [3]byte, initAddr byte) {
	slog.Info("New addr", "mac", mac[:], "addr", initAddr)
}

func (c *HwConfigurator) AddNode(mac [3]byte, initAddr byte) {
	reqAddr := initAddr
	newNode := true
	if reqAddr == 255 {
		reqAddr = c.firstFreeAddr()
	}
	for _, no := range c.nodes {
		if no.mac == mac {
			newNode = false
			break
		} else {
			if no.addr == reqAddr {
				reqAddr = c.firstFreeAddr()
			}
		}
	}
	if newNode {
		fmt.Println("New Node")
	} else {
		fmt.Println("Welcome back")
	}
	c.sendNodeNfo(mac, reqAddr)
	c.nodes = append(c.nodes, node{mac: mac, addr: reqAddr})

}
