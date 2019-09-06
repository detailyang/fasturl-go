package fasturl

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/stretchr/testify/require"
)

type T struct {
	haserr   bool
	input    string
	protocol string
	slashes  bool
	auth     string
	user     string
	pass     string
	host     string
	hostname string
	port     string
	path     string
	pathname string
	search   string
	query    string
	hash     string
	href     string
}

func TestFastURLPanic(t *testing.T) {
	for _, tt := range []string{
		"/#0/",
		"#0/",
		"#?",
		"?0",
		"0@",
	} {
		var f1 FastURL
		err := f1.Parse([]byte(tt))
		require.Nil(t, err)
		d1 := f1.Encode(nil)

		var f2 FastURL
		err = f2.Parse(d1)
		d2 := f2.Encode(nil)

		if !reflect.DeepEqual(d1, d2) {
			fmt.Printf("url1: %#v\n", string(d1))
			fmt.Printf("url2: %#v\n", string(d2))
			panic("fail")
		}

		require.Nil(t, err)
	}
}

func TestFastURLQuery(t *testing.T) {
	var f FastURL
	err := f.Parse([]byte("http://x.com?foo=bar&bar=1&bar=2"))
	require.Nil(t, err)
	v := f.GetQuery()

	z, ok := v.Get([]byte("foo"))
	require.True(t, ok)
	require.Equal(t, "bar", string(z))

	z, ok = v.Get([]byte("bar"))
	require.True(t, ok)
	require.Equal(t, "1", string(z))

	values := [][]byte{}
	v.GetAll([]byte("bar"), func(value []byte) bool {
		values = append(values, value)
		return true
	})
	require.True(t, ok)
	require.Equal(t, [][]byte{[]byte("1"), []byte("2")}, values)
}

func TestFastURLNormlizePathname(t *testing.T) {
	for _, tt := range []struct {
		Input  string
		Expect string
	}{
		{
			"/aa//bb",
			"/aa/bb",
		},
		{
			"/x///y/", "/x/y/",
		},
		{
			"/abc//de///fg////", "/abc/de/fg/",
		},
		{
			"/xxxx%2fyyy%2f%2F%2F", "/xxxx/yyy/",
		},
		{
			"/aaa/..", "/",
		},
		{
			"/aaa/bbb/ccc/../../ddd", "/aaa/ddd",
		},
		{
			"/a/b/../c/d/../e/..", "/a/c/",
		},
		{
			"/aaa/../../../../xxx", "/xxx",
		},
		{
			"/aaa%2Fbbb%2F%2E.%2Fxxx", "/aaa/xxx",
		},
		{
			"/a/./b/././c/./d.html", "/a/b/c/d.html",
		},
		{
			"./foo/", "/foo/",
		},
		{
			"./../.././../../aaa/bbb/../../../././../", "/",
		},
		{
			"./a/./.././../b/./foo.html", "/b/foo.html",
		},
	} {
		var d []byte
		d = NormalizePathname(d, []byte(tt.Input))
		require.Equal(t, tt.Expect, string(d), tt.Input)
	}
}

func TestFastURLParse(t *testing.T) {
	tt := []T{
		T{
			input:    "//some_path",
			host:     "some_path",
			pathname: "/",
		},
		T{
			input:  "http://\t",
			haserr: true,
		},
		T{
			input:  "http://a\r \t\n<b:b@c\r\nd/e?f",
			haserr: true,
		},
		T{
			input:    "HtTp://bing.com/search?q=dotnet#hash",
			protocol: "http",
			host:     "bing.com",
			pathname: "/search",
			query:    "q=dotnet",
			hash:     "#hash",
		},
		T{
			input:    "http://www.ExAmPlE.com/",
			protocol: "http",
			host:     "www.example.com",
			pathname: "/",
			query:    "",
			hash:     "",
		},
		T{
			input:    "http://user:pw@www.ExAmPlE.com/",
			protocol: "http",
			host:     "www.example.com",
			hostname: "www.example.com",
			port:     "",
			user:     "user",
			pass:     "pw",
			pathname: "/",
			query:    "",
			hash:     "",
		},
		T{
			input:    "git+ssh://git@github.com:npm/npm",
			protocol: "git+ssh",
			slashes:  true,
			auth:     "git",
			user:     "git",
			host:     "github.com:npm",
			port:     "",
			hostname: "github.com",
			pathname: "/npm",
		},
		T{
			input:    "http://a@b?@c",
			protocol: "http",
			auth:     "a",
			user:     "a",
			host:     "b",
			hostname: "b",
			pathname: "/",
			query:    "@c",
		},
		T{
			input:    "http://ليهمابتكلموشعربي؟.ي؟/",
			protocol: "http",
			host:     "ليهمابتكلموشعربي؟.ي؟",
			pathname: "/",
		},
		T{
			input:    "http://atpass:foo%40bar@127.0.0.1:8080/path?search=foo#bar",
			protocol: "http",
			host:     "127.0.0.1:8080",
			auth:     "atpass:foo@bar",
			user:     "atpass",
			pass:     "foo%40bar",
			hostname: "127.0.0.1",
			port:     "8080",
			pathname: "/path",
			query:    "search=foo",
			hash:     "#bar",
		},
		T{
			input:    "//user:pass@example.com:8000/foo/bar?baz=quux#frag",
			protocol: "",
			host:     "example.com:8000",
			auth:     "user:pass",
			user:     "user",
			pass:     "pass",
			port:     "8000",
			hostname: "example.com",
			hash:     "#frag",
			query:    "baz=quux",
			pathname: "/foo/bar",
		},
	}

	var f FastURL

	for i := range tt {
		f.Reset()
		err := Parse(&f, []byte(tt[i].input))
		input := tt[i].input
		if tt[i].haserr {
			require.NotNil(t, err, input)
			continue
		}
		require.Nil(t, err, input)
		require.Equal(t, tt[i].protocol, string(f.GetProtocol()), input)
		require.Equal(t, tt[i].host, string(f.GetHost()), input)
		require.Equal(t, tt[i].user, string(f.GetUser()), input)
		require.Equal(t, tt[i].pass, string(f.GetPass()), input)
		require.Equal(t, tt[i].pathname, string(f.GetPathname()), input)
		require.Equal(t, tt[i].hash, string(f.GetHash()), input)
	}
}
