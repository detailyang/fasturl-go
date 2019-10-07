package fasturl

import (
	"bytes"
	"unsafe"
)

func b2s(bs []byte) string {
	return *(*string)(unsafe.Pointer(&bs))
}

func s2b(s string) []byte {
	return *(*[]byte)(unsafe.Pointer(&s))
}

func toLowercsaeASCII(b []byte) {
	for i := range b {
		if b[i] >= 'A' && b[i] <= 'Z' {
			b[i] += 'a' - 'A'
		}
	}
}

func hasASCIIControl(b []byte) bool {
	for i := 0; i < len(b); i++ {
		if b[i] < ' ' || b[i] == 0x7f {
			return true
		}
	}
	return false
}

type encoding int

const (
	encodePath encoding = 1 + iota
	encodePathSegment
	encodeHost
	encodeZone
	encodeUserPassword
	encodeQueryComponent
	encodeFragment
)

func ishex(c byte) bool {
	switch {
	case '0' <= c && c <= '9':
		return true
	case 'a' <= c && c <= 'f':
		return true
	case 'A' <= c && c <= 'F':
		return true
	}
	return false
}

func unhex(c byte) byte {
	switch {
	case '0' <= c && c <= '9':
		return c - '0'
	case 'a' <= c && c <= 'f':
		return c - 'a' + 10
	case 'A' <= c && c <= 'F':
		return c - 'A' + 10
	}
	return 0
}

// UnescapeQuery unescapes query
//
// The original implement is from net/ur
func UnescapeQuery(dst, src []byte) ([]byte, error) {
	return unescape(dst, src, encodeQueryComponent)
}

// unescape unescapes a string; the mode specifies
// which section of the URL string is being unescaped.
//
// The original implement is from net/url
func unescape(dst, src []byte, mode encoding) ([]byte, error) {
	// Count %, check that they're well-formed.
	n := 0
	hasPlus := false
	for i := 0; i < len(src); {
		switch src[i] {
		case '%':
			n++
			if i+2 >= len(src) || !ishex(src[i+1]) || !ishex(src[i+2]) {
				src = src[i:]
				if len(src) > 3 {
					src = src[:3]
				}
				return nil, ErrFastURLInvalidCharacter
			}

			// Per https://tools.ietf.org/html/rfc3986#page-21
			// in the host component %-encoding can only be used
			// for non-ASCII bytes.
			// But https://tools.ietf.org/html/rfc6874#section-2
			// introduces %25 being allowed to escape a percent sign
			// in IPv6 scoped-address literals. Yay.
			if mode == encodeHost && unhex(src[i+1]) < 8 && string(src[i:i+3]) != "%25" {
				return nil, ErrFastURLInvalidCharacter
			}
			if mode == encodeZone {
				// RFC 6874 says basically "anything goes" for zone identifiers
				// and that even non-ASCII can be redundantly escaped,
				// but it seems prudent to restrict %-escaped bytes here to those
				// that are valid host name bytes in their unescaped form.
				// That is, you can use escaping in the zone identifier but not
				// to introduce bytes you couldn't just write directly.
				// But Windows puts spaces here! Yay.
				v := unhex(src[i+1])<<4 | unhex(src[i+2])
				if string(src[i:i+3]) != "%25" && v != ' ' && shouldEscape(v, encodeHost) {
					return nil, ErrFastURLInvalidCharacter
				}
			}
			i += 3
		case '+':
			hasPlus = mode == encodeQueryComponent
			i++
		default:
			if (mode == encodeHost || mode == encodeZone) && src[i] < 0x80 && shouldEscape(src[i], mode) {
				return nil, ErrFastURLInvalidCharacter
			}
			i++
		}
	}

	if n == 0 && !hasPlus {
		dst = append(dst, src...)
		return dst, nil
	}

	for i := 0; i < len(src); i++ {
		switch src[i] {
		case '%':
			dst = append(dst, unhex(src[i+1])<<4|unhex(src[i+2]))
			i += 2
		case '+':
			if mode == encodeQueryComponent {
				dst = append(dst, ' ')
			} else {
				dst = append(dst, '+')
			}
		default:
			dst = append(dst, src[i])
		}
	}
	return dst, nil
}

// EscapeQuery escapes the string so it can be safely placed
// inside a URL rawquery.
//
// The original implement is from net/url
func EscapeQuery(dst, src []byte) []byte {
	return escape(dst, src, encodeQueryComponent)
}

func escape(dst, src []byte, mode encoding) []byte {
	spaceCount, hexCount := 0, 0
	for i := 0; i < len(src); i++ {
		c := src[i]
		if shouldEscape(c, mode) {
			if c == ' ' && mode == encodeQueryComponent {
				spaceCount++
			} else {
				hexCount++
			}
		}
	}

	// Fastpath
	if spaceCount == 0 && hexCount == 0 {
		dst = append(dst, src...)
		return dst
	}

	for i := 0; i < len(src); i++ {
		switch c := src[i]; {
		case c == byte(' ') && mode == encodeQueryComponent:
			dst = append(dst, '+')
		case shouldEscape(c, mode):
			dst = append(dst, '%', "0123456789ABCDEF"[c>>4], "0123456789ABCDEF"[c&15])
		default:
			dst = append(dst, c)
		}
	}

	return dst
}

// Return true if the specified character should be escaped when
// appearing in a URL string, according to RFC 3986.
//
// Please be informed that for now shouldEscape does not check all
// reserved characters correctly. See golang.org/issue/5684.
//
// The original implement is from net/url
func shouldEscape(c byte, mode encoding) bool {
	// §2.3 Unreserved characters (alphanum)
	if 'a' <= c && c <= 'z' || 'A' <= c && c <= 'Z' || '0' <= c && c <= '9' {
		return false
	}

	if mode == encodeHost || mode == encodeZone {
		// §3.2.2 Host allows
		//	sub-delims = "!" / "$" / "&" / "'" / "(" / ")" / "*" / "+" / "," / ";" / "="
		// as part of reg-name.
		// We add : because we include :port as part of host.
		// We add [ ] because we include [ipv6]:port as part of host.
		// We add < > because they're the only characters left that
		// we could possibly allow, and Parse will reject them if we
		// escape them (because hosts can't use %-encoding for
		// ASCII bytes).
		switch c {
		case '!', '$', '&', '\'', '(', ')', '*', '+', ',', ';', '=', ':', '[', ']', '<', '>', '"':
			return false
		}
	}

	switch c {
	case '-', '_', '.', '~': // §2.3 Unreserved characters (mark)
		return false

	case '$', '&', '+', ',', '/', ':', ';', '=', '?', '@': // §2.2 Reserved characters (reserved)
		// Different sections of the URL allow a few of
		// the reserved characters to appear unescaped.
		switch mode {
		case encodePath: // §3.3
			// The RFC allows : @ & = + $ but saves / ; , for assigning
			// meaning to individual path segments. This package
			// only manipulates the path as a whole, so we allow those
			// last three as well. That leaves only ? to escape.
			return c == '?'

		case encodePathSegment: // §3.3
			// The RFC allows : @ & = + $ but saves / ; , for assigning
			// meaning to individual path segments.
			return c == '/' || c == ';' || c == ',' || c == '?'

		case encodeUserPassword: // §3.2.1
			// The RFC allows ';', ':', '&', '=', '+', '$', and ',' in
			// userinfo, so we must escape only '@', '/', and '?'.
			// The parsing of userinfo treats ':' as special so we must escape
			// that too.
			return c == '@' || c == '/' || c == '?' || c == ':'

		case encodeQueryComponent: // §3.4
			// The RFC reserves (so we must escape) everything.
			return true

		case encodeFragment: // §4.1
			// The RFC text is silent but the grammar allows
			// everything, so escape nothing.
			return false
		}
	}

	if mode == encodeFragment {
		// RFC 3986 §2.2 allows not escaping sub-delims. A subset of sub-delims are
		// included in reserved from RFC 2396 §2.2. The remaining sub-delims do not
		// need to be escaped. To minimize potential breakage, we apply two restrictions:
		// (1) we always escape sub-delims outside of the fragment, and (2) we always
		// escape single quote to avoid breaking callers that had previously assumed that
		// single quotes would be escaped. See issue #19917.
		switch c {
		case '!', '(', ')', '*':
			return false
		}
	}

	// Everything else must be escaped.
	return true
}

// NormalizePathname normalizes pathname
//
// The original implement is from fasthttp
func NormalizePathname(dst, src []byte) []byte {
	dst = dst[:0]
	if len(src) == 0 || src[0] != '/' {
		dst = append(dst, '/')
	}
	var err error
	dst, err = unescape(dst, src, encodePath)
	if err != nil {
		dst = append(dst, src...)
		return dst
	}

	// remove duplicate slashes
	b := dst
	bSize := len(b)
	for {
		n := bytes.Index(b, []byte("//"))
		if n < 0 {
			break
		}
		b = b[n:]
		copy(b, b[1:])
		b = b[:len(b)-1]
		bSize--
	}
	dst = dst[:bSize]

	// remove /./ parts
	b = dst
	for {
		n := bytes.Index(b, []byte("/./"))
		if n < 0 {
			break
		}
		nn := n + len([]byte("/./")) - 1
		copy(b[n:], b[nn:])
		b = b[:len(b)-nn+n]
	}

	// remove /foo/../ parts
	for {
		n := bytes.Index(b, []byte("/../"))
		if n < 0 {
			break
		}
		nn := bytes.LastIndexByte(b[:n], '/')
		if nn < 0 {
			nn = 0
		}
		n += len([]byte("/../")) - 1
		copy(b[nn:], b[n:])
		b = b[:len(b)-n+nn]
	}

	// remove trailing /foo/..
	n := bytes.LastIndex(b, []byte("/.."))
	if n >= 0 && n+len([]byte("/..")) == len(b) {
		nn := bytes.LastIndexByte(b[:n], '/')
		if nn < 0 {
			return []byte("/")
		}
		b = b[:nn+1]
	}

	return b
}
