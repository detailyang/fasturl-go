package fasturl

import "bytes"

type Path struct {
	path               []byte
	pathname           []byte
	normalizedPathname []byte
	parsequery         bool
	query              Query
	rawquery           []byte
	hash               []byte
}

// GetPathname gets the pathname
func (p *Path) GetPathname() []byte {
	return p.pathname
}

// SetPathname sets the pathname
func (p *Path) SetPathname(pn string) {
	p.pathname = append(p.pathname[:0], pn...)
}

// GetNormalizedPathname gets the normalized pathname
func (p *Path) GetNormalizedPathname() []byte {
	return p.normalizedPathname
}

// GetRawQuery gets the raw query string
func (p *Path) GetRawQuery() []byte {
	return p.rawquery
}

// GetQuery gets the query
func (p *Path) GetQuery() *Query {
	if !p.parsequery {
		p.query.Decode(p.rawquery)
		p.parsequery = true
	}
	return &p.query
}

// GetHash gets the hash
func (p *Path) GetHash() []byte {
	return p.hash
}

// SetHash sets the hash
func (p *Path) SetHash(hash string) {
	p.hash = append(p.hash[:0], hash...)
}

// Parse parse path to Path
func (p *Path) Parse(path []byte) error {
	return ParsePath(p, path)
}

// Reset resets the Path
func (p *Path) Reset() {
	p.pathname = p.pathname[:0]
	p.normalizedPathname = p.normalizedPathname[:0]
	p.rawquery = p.rawquery[:0]
	p.hash = p.hash[:0]
	p.parsequery = false
	p.query.Reset()
}

// ParsePath parse the path to Path
func ParsePath(f *Path, path []byte) error {
	// Find hash
	hashIndex := bytes.IndexByte(path, '#')
	if hashIndex >= 0 {
		f.hash = append(f.hash, path[hashIndex:]...)
		path = path[:hashIndex]
	}

	if hasASCIIControl(path) {
		return ErrFastURLInvalidCharacter
	}

	// Find query
	queryIndex := bytes.IndexByte(path, '?')
	if queryIndex >= 0 {
		f.rawquery = append(f.rawquery[:0], path[queryIndex+1:]...)
		path = path[:queryIndex]
	}

	pos := 0
	// Find path
	if len(path[pos:]) > 0 && path[pos] != '/' {
		f.pathname = append(f.pathname[:0], '/')
	}

	f.pathname = append(f.pathname, path[pos:]...)
	if len(f.pathname) == 0 {
		f.pathname = append(f.pathname, '/')
	}

	f.normalizedPathname = NormalizePathname(f.normalizedPathname[:0], f.pathname)

	return nil
}
