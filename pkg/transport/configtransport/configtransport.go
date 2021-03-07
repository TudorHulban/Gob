package configtransport

import (
	"strconv"
)

// TODO: load configuration from external source.

type Cfg struct {
	TerminationSecs uint
	Port            uint
	IP              string
	Protocol        string
}

const (
	TerminationSecs = 1
	Port            = 8080
	IP              = "localhost"
	Protocol        = "tcp"
	ACK             = "thank you"
)

func NewDefaultConfiguration() *Cfg {
	return &Cfg{
		TerminationSecs: TerminationSecs,
		Port:            Port,
		IP:              IP,
		Protocol:        Protocol,
	}
}

func (c *Cfg) Socket() string {
	return c.IP + ":" + strconv.FormatUint(uint64(c.Port), 10)
}
