package fasturl

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestQuery(t *testing.T) {
	t.Run("test set and get and del", func(t *testing.T) {
		var q Query

		q.SetBytes([]byte("a"), []byte("b"))
		d, ok := q.GetBytes([]byte("a"))
		require.Equal(t, true, ok)
		require.Equal(t, "b", string(d))

		q.SetBytes([]byte("a"), []byte("c"))
		d, ok = q.GetBytes([]byte("a"))
		require.Equal(t, true, ok)
		require.Equal(t, "c", string(d))

		q.DelBytes([]byte("a"))
		d, ok = q.GetBytes([]byte("a"))
		require.Equal(t, false, ok)
	})

	t.Run("test add and getall and delall", func(t *testing.T) {
		var q Query
		q.AddBytes([]byte("abcd"), []byte("defg"))
		q.AddBytes([]byte("abcd"), []byte("cccc"))
		q.AddBytes([]byte("abcd"), []byte("1234"))

		values := [][]byte{}
		q.GetAllBytes([]byte("abcd"), func(value []byte) bool {
			values = append(values, value)
			return true
		})
		require.Equal(t, [][]byte{[]byte("defg"), []byte("cccc"), []byte("1234")}, values)

		q.DelAllBytes([]byte("abcd"))
		values = [][]byte{}
		q.GetAllBytes([]byte("abcd"), func(value []byte) bool {
			values = append(values, value)
			return true
		})
		require.Equal(t, 0, len(values))
	})

	t.Run("test range and reset", func(t *testing.T) {
		var q Query
		q.AddBytes([]byte("abcd"), []byte("defg"))
		q.AddBytes([]byte("abcd"), []byte("cccc"))
		q.AddBytes([]byte("abcd"), []byte("1234"))
		values := [][]byte{}
		q.Range(func(key, value []byte) bool {
			values = append(values, value)
			return true
		})
		require.Equal(t, [][]byte{[]byte("defg"), []byte("cccc"), []byte("1234")}, values)
		q.Reset()
		values = [][]byte{}
		q.Range(func(key, value []byte) bool {
			values = append(values, value)
			return true
		})
		require.Equal(t, 0, len(values))
	})

	t.Run("test decode and encode", func(t *testing.T) {
		var q Query
		input := []byte("a=1&b=2&c=3&a=1")
		err := q.Decode(input)
		require.Nil(t, err)

		d, ok := q.GetBytes([]byte("a"))
		require.Equal(t, true, ok)
		require.Equal(t, "1", string(d))

		d, ok = q.GetBytes([]byte("b"))
		require.Equal(t, true, ok)
		require.Equal(t, "2", string(d))

		d, ok = q.GetBytes([]byte("c"))
		require.Equal(t, true, ok)
		require.Equal(t, "3", string(d))

		values := [][]byte{}
		q.GetAllBytes([]byte("a"), func(value []byte) bool {
			values = append(values, value)
			return true
		})
		require.Equal(t, [][]byte{[]byte("1"), []byte("1")}, values)

		var dd []byte
		dd = q.Encode(dd)
		require.Equal(t, string(input), string(dd))
	})
}
