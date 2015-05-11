package urit

import (
	"strings"
)

type Variables map[string]string

type URI string

// Expand expands a uri removing any non-supplied variables
func (u URI) Expand(v Variables) URI {
	return u.expand(true, v)
}

// Inspect expand a uri including non-supplied variables (as expressions) back
// into the returned uri. It attemps to maintain the order in which the
// variables where original defined.
func (u URI) Inspect(v Variables) URI {
	return u.expand(false, v)
}

func (u URI) String(v Variables) string {
	return string(u.Expand(v))
}

// expand splits a URI into it's multiple expression pieces and expands each one
// then replacing them back in the original uri string
func (u URI) expand(d bool, v Variables) URI {
	s := string(u)
	for _, e := range Split(u) {
		s = strings.Replace(s, e.Match, e.Expand(d, v), -1)
	}

	return URI(s)
}
