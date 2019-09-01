package fasturl

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestQuery(t *testing.T) {
	t.Run("test set and get and del", func(t *testing.T) {
		var q Query

		q.Set([]byte("a"), []byte("b"))
		d, ok := q.Get([]byte("a"))
		require.Equal(t, true, ok)
		require.Equal(t, "b", string(d))

		q.Set([]byte("a"), []byte("c"))
		d, ok = q.Get([]byte("a"))
		require.Equal(t, true, ok)
		require.Equal(t, "c", string(d))

		q.Del([]byte("a"))
		d, ok = q.Get([]byte("a"))
		require.Equal(t, false, ok)
	})

	t.Run("test add and getall and delall", func(t *testing.T) {
		var q Query
		q.Add([]byte("abcd"), []byte("defg"))
		q.Add([]byte("abcd"), []byte("cccc"))
		q.Add([]byte("abcd"), []byte("1234"))

		values := [][]byte{}
		q.GetAll([]byte("abcd"), func(value []byte) bool {
			values = append(values, value)
			return true
		})
		require.Equal(t, [][]byte{[]byte("defg"), []byte("cccc"), []byte("1234")}, values)

		q.DelAll([]byte("abcd"))
		values = [][]byte{}
		q.GetAll([]byte("abcd"), func(value []byte) bool {
			values = append(values, value)
			return true
		})
		require.Equal(t, 0, len(values))
	})

	t.Run("test range and reset", func(t *testing.T) {
		var q Query
		q.Add([]byte("abcd"), []byte("defg"))
		q.Add([]byte("abcd"), []byte("cccc"))
		q.Add([]byte("abcd"), []byte("1234"))
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

		d, ok := q.Get([]byte("a"))
		require.Equal(t, true, ok)
		require.Equal(t, "1", string(d))

		d, ok = q.Get([]byte("b"))
		require.Equal(t, true, ok)
		require.Equal(t, "2", string(d))

		d, ok = q.Get([]byte("c"))
		require.Equal(t, true, ok)
		require.Equal(t, "3", string(d))

		values := [][]byte{}
		q.GetAll([]byte("a"), func(value []byte) bool {
			values = append(values, value)
			return true
		})
		require.Equal(t, [][]byte{[]byte("1"), []byte("1")}, values)

		var dd []byte
		dd = q.Encode(dd)
		require.Equal(t, string(input), string(dd))
	})
}
