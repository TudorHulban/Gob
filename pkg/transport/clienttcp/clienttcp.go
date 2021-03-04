package clienttcp

import (
	"bufio"
	"main/pkg/transport/config"
	"net"
)

type ClientTCP struct {
	configtransport.Cfg
}

// NewClient Constructor defined for convenience.
func NewClient(cfg *configtransport.Cfg) (*ClientTCP, error) {
	result := new(ClientTCP)

	if cfg == nil {
		result.Cfg = *configtransport.NewDefaultConfiguration()
		return result, nil
	}

	result.Cfg = *cfg
	return result, nil
}

func (c *ClientTCP) Send(payload []byte) (string, error) {
	conn, errDial := net.Dial(c.Protocol, c.Cfg.Socket())
	if errDial != nil {
		return "", errDial
	}

	return bufio.NewReader(conn).ReadString('\n')
}
