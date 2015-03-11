package urit

import (
	"path/filepath"
	"strings"
)

type Operator string

// Operator : [Prefix, Joiner]
var opmap = map[Operator][]string{
	"/": {"/", "/"},
	"#": {"#", ","},
	"?": {"?", "&"},
	"&": {"&", "&"},
	";": {";", ";"},
	".": {".", "."},
}

func (o Operator) Join(a []string) (string, error) {
	var s string

	switch o {
	case "/":
		s = filepath.Join(append([]string{"/"}, a...)...)

	case "#", "?", ".", ";", "&":
		m := opmap[o]
		s = m[0] + strings.Join(a, m[1])

	default:
		s = strings.Join(a, ",")
	}

	return s, nil
}
