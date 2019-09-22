package fasturl

import (
	"bytes"
	"errors"
)

var (
	// ErrFastURLInvalidCharacter indicates the url contains invalid character.
	ErrFastURLInvalidCharacter = errors.New("fasturl: invalid character in url")
)

// FastURL presents the URL.
//
// ┌─────────────────────────────────────────────────────────────────────────────────────────────-───┐
// │                                              href                                               │
// ├──────────┬-──┬─────────────────────┬────────────────────────┬───────────────────────────┬───────┤
// │ protocol │   │        auth         │          host          │           path            │ hash  │
// │          │   │                     ├─────────────────┬──────┼──────────┬────────────────┤       │
// │          │   │                     │    hostname     │ port │ pathname │     search     │       │
// │          │   │                     │                 │      │          ├─┬──────────────┤       │
// │          │   │                     │                 │      │          │ │    query     │       │
// "  https    ://    user   :   pass   @ sub.example.com : 8080   /p/a/t/h  ?  query=string   #hash "
// │          │   │          │          │    hostname     │ port │          │                │       │
// │          │   │          │          ├─────────────────┴──────┤          │                │       │
// │ protocol │   │ username │ password │          host          │          │                │       │
// ├──────────┴-──┼──────────┴──────────┼────────────────────────┤          │                │       │
// │   origin     │                     │         origin         │ pathname │     search     │ hash  │
// ├───────────-──┴─────────────────────┴────────────────────────┴──────────┴────────────────┴───────┤
// │                                              href                                               │
// └──────────────────────────────────────────────────────────────────────────────────────────-──────┘
// (all spaces in the "" line should be ignored — they are purely for formatting)
// More detail see https://url.spec.whatwg.org/
//
//[protocol:][//[auth@]host][/]pathname[?query][#hash]
type FastURL struct {
	protocol           []byte
	auth               []byte
	user               []byte
	pass               []byte
	host               []byte
	hostname           []byte
	port               []byte
	path               []byte
	pathname           []byte
	normalizedPathname []byte
	parsequery         bool
	query              Query
	rawquery           []byte
	hash               []byte
}

// GetProtocol gets the protocol.
func (f *FastURL) GetProtocol() []byte {
	return f.protocol
}

// SetProtocol sets the protocol.
func (f *FastURL) SetProtocol(p string) {
	f.protocol = append(f.protocol[:0], p...)
}

// GetAuth gets the auth.
func (f *FastURL) GetAuth() []byte {
	return f.auth
}

// GetUser gets the username.
func (f *FastURL) GetUser() []byte {
	return f.user
}

// SetUser sets the username.
func (f *FastURL) SetUser(username string) {
	f.user = append(f.user[:0], username...)
}

// GetPass gets the password.
func (f *FastURL) GetPass() []byte {
	return f.pass
}

// SetPass sets the password.
func (f *FastURL) SetPass(password string) {
	f.pass = append(f.pass[:0], password...)
}

// GetHost gets the host which include hostname and port.
func (f *FastURL) GetHost() []byte {
	return f.host
}

// GetHostname gets the hostname.
func (f *FastURL) GetHostname() []byte {
	return f.hostname
}

// SetHostname sets the hostname.
func (f *FastURL) SetHostname(hostname string) {
	f.hostname = append(f.hostname[:0], hostname...)
}

// GetPort gets the port.
func (f *FastURL) GetPort() []byte {
	return f.port
}

// SetPort sets the port.
func (f *FastURL) SetPort(port string) {
	f.port = append(f.port[:0], port...)
}

// GetPathname gets the pathname.
func (f *FastURL) GetPathname() []byte {
	return f.pathname
}

// SetPathname sets the pathname.
func (f *FastURL) SetPathname(p string) {
	f.pathname = append(f.pathname[:0], p...)
}

// GetNormalizedPathname gets the normalized pathname.
func (f *FastURL) GetNormalizedPathname() []byte {
	return f.normalizedPathname
}

// GetRawQuery gets the raw query string.
func (f *FastURL) GetRawQuery() []byte {
	return f.rawquery
}

// GetQuery gets the query.
func (f *FastURL) GetQuery() *Query {
	if !f.parsequery {
		f.query.Decode(f.rawquery)
		f.parsequery = true
	}
	return &f.query
}

// GetHash gets the hash.
func (f *FastURL) GetHash() []byte {
	return f.hash
}

// SetHash sets the hash.
func (f *FastURL) SetHash(hash string) {
	f.hash = append(f.hash[:0], hash...)
}

// Parse parses the url.
func (f *FastURL) Parse(url []byte) error {
	return Parse(f, url)
}

// Encode encodes to []byte.
func (f *FastURL) Encode(b []byte) []byte {
	if len(f.protocol) > 0 {
		b = append(b, f.protocol...)
		b = append(b, ':')
	}

	var appendslash bool
	if len(f.user) > 0 || len(f.pass) > 0 {
		b = append(b, "//"...)
		appendslash = true
		b = append(b, f.user...)
		b = append(b, ':')
		b = append(b, f.pass...)
		b = append(b, '@')
	}

	if len(f.hostname) > 0 {
		if !appendslash {
			b = append(b, "//"...)
			appendslash = true
		}
		b = append(b, f.hostname...)
	}

	if len(f.port) > 0 {
		if !appendslash {
			b = append(b, "//"...)
			appendslash = true
		}
		b = append(b, ':')
		b = append(b, f.port...)
	}

	if len(f.pathname) > 0 {
		b = append(b, f.pathname...)
	}

	if f.query.Len() > 0 {
		b = append(b, '?')
		b = f.query.Encode(b)
	}

	if len(f.hash) > 0 {
		b = append(b, f.hash...)
	}

	return b
}

// Reset resets the FastURL.
func (f *FastURL) Reset() {
	f.protocol = f.protocol[:0]
	f.auth = f.auth[:0]
	f.user = f.user[:0]
	f.pass = f.pass[:0]
	f.host = f.host[:0]
	f.hostname = f.hostname[:0]
	f.port = f.port[:0]
	f.pathname = f.pathname[:0]
	f.normalizedPathname = f.normalizedPathname[:0]
	f.rawquery = f.rawquery[:0]
	f.hash = f.hash[:0]
	f.parsequery = false
	f.query.Reset()
}

// ParseWithoutProtocol parses the url to FastURL without protocol.
func ParseWithoutProtocol(f *FastURL, url []byte) error {
	return parse(f, url, parseOption{
		ParseProtocol: false,
	})
}

// Parse parses the url to FastURL.
func Parse(f *FastURL, url []byte) error {
	return parse(f, url, parseOption{
		ParseProtocol: true,
	})
}

type parseOption struct {
	ParseProtocol bool
}

func parse(f *FastURL, url []byte, o parseOption) error {
	// Find hash
	hashIndex := bytes.IndexByte(url, '#')
	if hashIndex >= 0 {
		f.hash = append(f.hash, url[hashIndex:]...)
		url = url[:hashIndex]
	}

	if hasASCIIControl(url) {
		return ErrFastURLInvalidCharacter
	}

	// Find query
	queryIndex := bytes.IndexByte(url, '?')
	if queryIndex >= 0 {
		f.rawquery = append(f.rawquery[:0], url[queryIndex+1:]...)
		url = url[:queryIndex]
	}

	pos := 0
	if o.ParseProtocol {
		// Trim //
		if len(url) >= 2 && string(url[0:2]) == "//" {
			pos += 2
		} else {
			// Find protocol
			pi := bytes.IndexByte(url[pos:], ':')
			if pi >= 0 {
				f.protocol = append(f.protocol[:0], url[:pi]...)
				toLowercsaeASCII(f.protocol)
				pos += pi + 1
			}
			// Trim //
			if len(url[pos:]) >= 2 && string(url[pos:pos+2]) == "//" {
				pos += 2

			}
		}
	}

	// Find auth
	// TODO(detailylang): follow https://url.spec.whatwg.org/#authority-state
	ai := bytes.IndexByte(url[pos:], '@')
	if ai >= 0 {
		f.auth = append(f.auth[:0], url[pos:pos+ai]...)
		// Find :
		ci := bytes.IndexByte(f.auth, ':')
		if ci >= 0 {
			f.user = append(f.user[:0], f.auth[:ci]...)
			f.pass = append(f.pass[:0], f.auth[ci+1:]...)
		} else {
			f.user = append(f.user[:0], f.auth...)
		}

		pos += ai + 1
	}

	// Find host
	hi := bytes.IndexByte(url[pos:], '/')
	if hi >= 0 {
		f.host = append(f.host[:0], url[pos:pos+hi]...)
		pos += hi

	} else {
		f.host = append(f.host[:0], url[pos:]...)
		pos = len(url)
	}
	toLowercsaeASCII(f.host)

	if len(f.host) > 0 {
		// find :
		ci := bytes.IndexByte(f.host, ':')
		if ci >= 0 {
			f.hostname = append(f.hostname[:0], f.host[:ci]...)
			f.port = append(f.port[:0], f.host[ci+1:]...)
		} else {
			f.hostname = append(f.hostname[:0], f.host...)
		}
	}

	// Find path
	if len(url[pos:]) > 0 && url[pos] != '/' {
		f.pathname = append(f.pathname[:0], '/')
	}

	f.pathname = append(f.pathname, url[pos:]...)
	if len(f.pathname) == 0 {
		f.pathname = append(f.pathname, '/')
	}

	f.normalizedPathname = NormalizePathname(f.normalizedPathname[:0], f.pathname)

	return nil
}
