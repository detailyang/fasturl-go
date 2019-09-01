<p align="center">
  <b>
    <span style="font-size:larger;">fasturl-go</span>
  </b>
  <br />
   <a href="https://travis-ci.org/detailyang/fasturl-go"><img src="https://travis-ci.org/detailyang/fasturl-go.svg?branch=master" /></a>
   <a href="https://ci.appveyor.com/project/detailyang/fasturl-go"><img src="https://ci.appveyor.com/api/projects/status/hbpj944ankoy9sh5?svg=true" /></a>
   <br />
   <b>fasturl-go is a yet another url parser but zero allocted which is more faster then net/url</b>
</p>

````bash
go test -v -benchmem -run="^$" github.com/detailyang/fasturl-go/fasturl -bench Benchmark
goos: darwin
goarch: amd64
pkg: github.com/detailyang/fasturl-go/fasturl
BenchmarkFastURLParseCase1-8   	 5639812	       190 ns/op	       0 B/op	       0 allocs/op
BenchmarkNetURLParseCase1-8    	 1459245	       824 ns/op	     352 B/op	       7 allocs/op
BenchmarkFastURLParseCase2-8   	 3414703	       364 ns/op	       0 B/op	       0 allocs/op
BenchmarkNetURLParseCase2-8    	  784420	      1461 ns/op	     528 B/op	       8 allocs/op
BenchmarkQuery-8               	 6587886	       170 ns/op	       0 B/op	       0 allocs/op
PASS
ok  	github.com/detailyang/fasturl-go/fasturl	7.405s
````
