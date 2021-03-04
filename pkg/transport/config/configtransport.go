package configtransport

import "strconv"

// TODO: load configuration from external source.

type Cfg struct {
	Port     uint
	IP       string
	Protocol string
}

const (
	Port     = 8080
	IP       = "localhost"
	Protocol = "tcp"
)

func NewDefaultConfiguration() *Cfg {
	return &Cfg{
		Port:     Port,
		IP:       IP,
		Protocol: Protocol,
	}
}

func (c *Cfg) Socket() string {
	return c.IP + ":" + strconv.FormatUint(uint64(c.Port), 10)
}
