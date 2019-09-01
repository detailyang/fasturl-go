package fasturl

import (
	"net/url"
	"testing"
)

func BenchmarkFastURLParseCase1(b *testing.B) {
	benchmarkFastURLParse(b, []byte("http://www.example.com/example?aaaa=1#bc"))
}

func BenchmarkNetURLParseCase1(b *testing.B) {
	benchmarkNetURLParse(b, []byte("http://www.example.com/example?aaaa=1#bc"))
}

func BenchmarkFastURLParseCase2(b *testing.B) {
	benchmarkFastURLParse(b, []byte("http://www.example.com/example/asdfasdf/../asdfasdf/1234?aaaa=1&q=2&aaa=3&fff=4&cc=3&123=f&ccc=3#bc"))
}

func BenchmarkNetURLParseCase2(b *testing.B) {
	benchmarkNetURLParse(b, []byte("http://www.example.com/example/asdfasdf/../asdfasdf/1234?aaaa=1&q=2&aaa=3&fff=4&cc=3&123=f&ccc=3#bc"))
}

func benchmarkNetURLParse(b *testing.B, input []byte) {
	var u url.URL
	uu := string(input)
	for i := 0; i < b.N; i++ {
		_, err := u.Parse(uu)
		if err != nil {
			b.Fatal("unexpected error", err)
		}
	}
}

func benchmarkFastURLParse(b *testing.B, url []byte) {
	var f FastURL
	for i := 0; i < b.N; i++ {
		f.Reset()
		err := f.Parse(url)
		if err != nil {
			b.Fatal("unexpected error", err)
		}
	}
}
