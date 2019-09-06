package fasturl

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestPath(t *testing.T) {
	var p Path

	x := []byte("/abcdef?d=1#hash")
	err := p.Parse(x)
	require.Nil(t, err)
	require.Equal(t, "/abcdef", string(p.GetPathname()))
	require.Equal(t, "#hash", string(p.GetHash()))
	v, ok := p.GetQuery().Get([]byte("d"))
	require.Equal(t, true, ok)
	require.Equal(t, "1", string(v))
}
