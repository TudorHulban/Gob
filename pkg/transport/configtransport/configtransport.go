package configtransport

import (
	"os"
	"strconv"

	"github.com/TudorHulban/log"
)

// TODO: load configuration from external source.

type Cfg struct {
	Port     uint
	IP       string
	Protocol string
	L        *log.LogInfo
}

const (
	Port     = 8080
	IP       = "localhost"
	Protocol = "tcp"
	ACK      = "thank you"
)

func NewDefaultConfiguration() *Cfg {
	return &Cfg{
		Port:     Port,
		IP:       IP,
		Protocol: Protocol,
		L:        log.New(log.DEBUG, os.Stdout, true),
	}
}

func (c *Cfg) Socket() string {
	return c.IP + ":" + strconv.FormatUint(uint64(c.Port), 10)
}
