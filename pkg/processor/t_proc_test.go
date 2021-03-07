package processor

import (
	"main/pkg/seri/serigob"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestProcGob(t *testing.T) {
	var g serigob.MGob

	p, errProc := NewProc(GobProcessing(g))
	require.Nil(t, errProc)
	require.Equal(t, p.Type, "Gob")
}
