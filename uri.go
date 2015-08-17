package urit

import (
	"strings"
)

type Variables map[string]string

type URI string

// Expand expands a uri removing any non-supplied variables
func (u URI) Expand(vars Variables) URI {
	return u.expand(true, vars)
}

// Inspect expand a uri including non-supplied variables (as expressions) back
// into the returned uri. It attemps to maintain the order in which the
// variables where original defined.
func (u URI) Inspect(vars Variables) URI {
	return u.expand(false, vars)
}

func (u URI) String(vars Variables) string {
	return string(u.Expand(vars))
}

// expand splits a URI into it's multiple expression pieces and expands each one
// then replacing them back in the original uri string
func (u URI) expand(pact bool, vars Variables) URI {
	str := string(u)
	for _, v := range Split(u) {
		str = strings.Replace(str, v.Match, v.Expand(pact, vars), -1)
	}

	return URI(str)
}
