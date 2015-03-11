package urit

import (
	"path/filepath"
	"strings"
)

type Operator string

// opmap => Operator: [Prefix, Joiner]
var opmap = map[Operator][]string{
	// op-level2
	"+": {"", ","},
	"#": {"#", ","},

	// op-level3
	".": {".", "."},
	"/": {"/", "/"},
	";": {";", ";"},
	"?": {"?", "&"},
	"&": {"&", "&"},
}

// Join joins an array of strings based on the format of the given operator
func (o Operator) Join(a []string) string {
	var s string

	switch o {
	case "/":
		s = filepath.Join(append([]string{"/"}, a...)...)

	case "#", "?", ".", ";", "&":
		m := opmap[o]
		s = m[0] + strings.Join(a, m[1])

	default: // handles +, `{+var}` and blank, `{var}`
		s = strings.Join(a, ",")
	}

	return s
}
