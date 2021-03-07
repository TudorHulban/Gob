package clienttcp

import (
	"bufio"
	"main/pkg/seri"
	"main/pkg/transport/configtransport"
	"net"
)

type ClientTCP struct {
	Payload []byte

	*configtransport.Cfg
}

// NewClient Constructor defined for convenience.
func NewClient(cfg *configtransport.Cfg) (*ClientTCP, error) {
	result := new(ClientTCP)

	if cfg == nil {
		c, err := configtransport.NewDefaultConfiguration()
		if err != nil {
			return nil, err
		}

		result.Cfg = c
		return result, nil
	}

	result.Cfg = cfg
	return result, nil
}

func (c *ClientTCP) PreprocessMsg(m seri.Message) *ClientTCP {
	var err error
	c.Payload, err = c.P.MsgEncode(m)

	c.L.Print("Error:", err)
	return c
}

func (c *ClientTCP) Send() (string, error) {
	c.L.Printf("Sending message to socket %s.", c.Socket())

	conn, errDial := net.Dial(c.Protocol, c.Cfg.Socket())
	if errDial != nil {
		c.L.Debug(errDial)
		return "", errDial
	}

	_, errSend := conn.Write(c.Payload)
	if errSend != nil {
		c.L.Debug(errSend)
		return "", errSend
	}

	return bufio.NewReader(conn).ReadString('\n')
}
