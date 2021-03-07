package processor

// Package processor added for structuring the processing with the help of serialization interface.

import (
	"io"
	"main/pkg/seri"
	"main/pkg/seri/serigob"

	"github.com/TudorHulban/log"
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

func Logger(level int, writeTo io.Writer) Option {
	return func(p *Proc) error {
		p.L = log.New(level, writeTo, true)
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
		return errDec
	}

	// TODO: add logic based on message
	p.L.Print("Message:", msg)

	return nil
}

func (p *Proc) MsgEncode(msg seri.Message) ([]byte, error) {
	return p.Actions.Encode(msg)
}
