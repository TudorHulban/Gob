package configtransport

import (
	"main/pkg/processor"
	"main/pkg/seri/serigob"
	"os"
	"strconv"

	"github.com/TudorHulban/log"
	"github.com/pkg/errors"
)

// TODO: load configuration from external source.

type Cfg struct {
	TerminationSecs uint
	Port            uint
	IP              string
	Protocol        string
	P               *processor.Proc
	L               *log.LogInfo
}

const (
	TerminationSecs = 9
	Port            = 8080
	IP              = "localhost"
	Protocol        = "tcp"
	ACK             = "thank you"
)

func NewDefaultConfiguration() (*Cfg, error) {
	var g serigob.MGob

	p, err := processor.NewProc(processor.GobProcessing(g))
	if err != nil {
		return nil, errors.WithMessage(err, "could not create default transport configuration")
	}

	return &Cfg{
		TerminationSecs: TerminationSecs,
		Port:            Port,
		IP:              IP,
		Protocol:        Protocol,
		P:               p,
		L:               log.New(log.DEBUG, os.Stdout, true),
	}, nil
}

func (c *Cfg) Socket() string {
	return c.IP + ":" + strconv.FormatUint(uint64(c.Port), 10)
}
