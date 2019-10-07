package fasturl

import "bytes"

// QueryPair holds the name & value
type QueryPair struct {
	name  []byte
	value []byte
}

func (q *QueryPair) set(name, value []byte) {
	q.name = append(q.name[:0], name...)
	q.value = append(q.value[:0], value...)
}

func (q *QueryPair) reset() {
	q.name = q.name[:0]
	q.value = q.value[:0]
}

// Query holds slice of QueryPair
type Query struct {
	pairs []QueryPair
}

// Len returns the length of query
func (q *Query) Len() int {
	return len(q.pairs)
}

// Reset resets the query
func (q *Query) Reset() {
	for i := range q.pairs {
		q.pairs[i].reset()
	}
	q.pairs = q.pairs[:0]
}

func (q *Query) GetAll(name string, fn func(value []byte) bool) {
	q.GetAllBytes(s2b(name), fn)
}

// GetAllBytes returns all query who name is equal to name
func (q *Query) GetAllBytes(name []byte, fn func(value []byte) bool) {
	for i := range q.pairs {
		pair := q.pairs[i]
		if bytes.Equal(pair.name, name) {
			if !fn(pair.value) {
				return
			}
		}
	}
}

func (q *Query) Get(name string) ([]byte, bool) {
	return q.GetBytes(s2b(name))
}

// GetBytes gets the value from name
func (q *Query) GetBytes(name []byte) ([]byte, bool) {
	for i := range q.pairs {
		pair := q.pairs[i]
		if bytes.Equal(pair.name, name) {
			return pair.value, true
		}
	}
	return nil, false
}

func (q *Query) Del(name string) {
	q.DelBytes(s2b(name))
}

// DelBytes dels the name
func (q *Query) DelBytes(name []byte) {
	q.del(name, false)
}

func (q *Query) DelAll(name string) {
	q.DelAllBytes(s2b(name))
}

// DelAllBytes dels all the value who name is equal to name
func (q *Query) DelAllBytes(name []byte) {
	q.del(name, true)
}

func (q *Query) del(name []byte, all bool) {
	for i := 0; i < len(q.pairs); i++ {
		pair := q.pairs[i]
		if bytes.Equal(pair.name, name) {
			q.pairs = append(q.pairs[:i], q.pairs[i+1:]...)
			if !all {
				return
			}
			i--
		}
	}
}

func (q *Query) Add(name, value string) {
	q.AddBytes(s2b(name), s2b(value))
}

// AddBytes adds the value
func (q *Query) AddBytes(name, value []byte) {
	pair := q.alloc()
	pair.set(name, value)
}

func (q *Query) Set(name, value string) {
	q.SetBytes(s2b(name), s2b(value))
}

// SetBytes sets or adds the name
func (q *Query) SetBytes(name, value []byte) {
	for i := range q.pairs {
		if bytes.Equal(q.pairs[i].name, name) {
			q.pairs[i].value = append(q.pairs[i].value[:0], value...)
			return
		}
	}

	q.AddBytes(name, value)
}

func (q *Query) alloc() *QueryPair {
	n := len(q.pairs)
	c := cap(q.pairs)
	if n == c {
		q.pairs = append(q.pairs, make([]QueryPair, 4)...)
	}
	q.pairs = q.pairs[:n+1]
	return &q.pairs[n]
}

func (q *Query) removeLastPair() {
	if n := len(q.pairs); n >= 1 {
		q.pairs = q.pairs[:n-1]
	}
}

// Range ranges the name
func (q *Query) Range(fn func(name, value []byte) bool) {
	for i := range q.pairs {
		if !fn(q.pairs[i].name, q.pairs[i].value) {
			return
		}
	}
}

// Decode decodes the querystring
func (q *Query) Decode(b []byte) error {
	return ParseQuery(q, b)
}

// Encode encodes the query to []byte
func (q *Query) Encode(b []byte) []byte {
	for i := range q.pairs {
		p := q.pairs[i]
		b = EscapeQuery(b, p.name)
		b = append(b, '=')
		b = EscapeQuery(b, p.value)
		if i < len(q.pairs)-1 {
			b = append(b, '&')
		}
	}
	return b
}

// ParseQuery parses the querystring to query
func ParseQuery(q *Query, query []byte) error {
	var err error
	for len(query) > 0 {
		key := query
		if i := bytes.IndexAny(key, "&;"); i >= 0 {
			key, query = key[:i], key[i+1:]
		} else {
			query = nil
		}

		if len(key) == 0 {
			continue
		}

		var value []byte
		if i := bytes.IndexByte(key, '='); i >= 0 {
			key, value = key[:i], key[i+1:]
		}

		pair := q.alloc()

		pair.name, err = UnescapeQuery(pair.name, key)
		if err != nil {
			q.removeLastPair()
			return err
		}

		pair.value, err = UnescapeQuery(pair.value, value)
		if err != nil {
			q.removeLastPair()
			return err
		}
	}

	return err
}
