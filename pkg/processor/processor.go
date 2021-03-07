package processor

// Package processor added for structuring the processing with the help of serialization interface.

import (
	"main/pkg/seri"
	"main/pkg/seri/serigob"

	"github.com/TudorHulban/log"
	"github.com/pkg/errors"
)

type Proc struct {
	Type    string
	Actions seri.ISerialization
	L       *log.LogInfo
}

// Option is functional option for processor configuration.
type Option func(p *Proc) error

// GobProcessing adds Gob processing to processor.
func GobProcessing(g serigob.MGob) Option {
	return func(p *Proc) error {
		p.Type = "Gob"
		p.Actions = g
		return nil
	}
}

func Logger(l *log.LogInfo) Option {
	return func(p *Proc) error {
		p.L = l
		return nil
	}
}

// NewProc Constructor using functional options for more detailed configuration.
func NewProc(opts ...Option) (*Proc, error) {
	p := new(Proc)

	for _, o := range opts {
		if err := o(p); err != nil {
			return nil, err
		}
	}

	if p.L == nil {
		p.L = log.New(log.NADA, nil, false)
	}

	return p, nil
}

// MsgDecode Method should have the capability to decode based on decoding function that
// was injected in constructor type.
func (p *Proc) MsgDecode(payload []byte) error {
	msg, errDec := p.Actions.Decode(payload)
	if errDec != nil {
		return errors.WithMessage(errDec, "could not decode message")
	}

	// TODO: add logic based on message
	p.L.Print("Message is:", msg)

	return nil
}

func (p *Proc) MsgEncode(msg seri.Message) ([]byte, error) {
	return p.Actions.Encode(msg)
}
