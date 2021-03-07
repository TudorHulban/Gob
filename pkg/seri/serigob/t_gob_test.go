package serigob

import (
	"main/pkg/seri"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestOps(t *testing.T) {
	m := seri.Message{
		ID:      1,
		Payload: "xxx",
	}

	var ops MGob
	e, errEnc := ops.Encode(m)

	require.Nil(t, errEnc)
	require.NotNil(t, e, "encoded value should not be nil")

	d, errDec := ops.Decode(e)

	require.Nil(t, errDec)
	require.Equal(t, m, d)
}
