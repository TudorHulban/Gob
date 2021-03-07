package processor

import (
	"main/pkg/seri"
	"main/pkg/seri/serigob"
)

type Proc struct {
	Type    string
	Actions seri.ISerialization
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

// NewProc Constructor using functional options for more detailed configuration.
// Only one option up to now though.
func NewProc(opts ...Option) (*Proc, error) {
	p := new(Proc)

	for _, o := range opts {
		if err := o(p); err != nil {
			return nil, err
		}

	return p, nil
}
