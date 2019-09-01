package fasturl

import "testing"

func BenchmarkQuery(b *testing.B) {
	url := []byte("aaaa=1&b=2&c=3&g=3")

	var q Query
	for i := 0; i < b.N; i++ {
		q.Reset()
		err := q.Decode(url)
		if err != nil {
			b.Fatal("unexpected error", err)
		}
	}
}
