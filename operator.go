package urit

import (
	"fmt"
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

	// op for continuation, not part of the original proposal. Used to handle
	// continuations on + and # operators use cases
	",": {",", ","},
}

// Join joins an array of strings based on the format of the given operator
func (o Operator) Join(a []string) string {
	var s string

	switch o {
	case "/":
		for i, v := range a {
			if v[:1] != "{" {
				a[i] = "/" + v
			}
		}
		s = filepath.Join(strings.Join(a, ""))

	case "#", "?", ".", ";", "&":
		m := opmap[o]
		for i, v := range a {
			if i > 0 && v[:1] != "{" {
				a[i] = m[1] + v
			}
		}
		s = m[0] + strings.Join(a, "")

	default: // handles +, `{+var}` and blank, `{var}`
		s = strings.Join(a, ",")
	}

	return s
}

// left and right expression delimiters { }
const (
	dl = "{"
	dr = "}"
)

func (o Operator) NewExpression(a []string) string {
	if len(a) == 0 {
		return ""
	}

	m, ok := opmap[o]
	if !ok {
		m = []string{"", ","}
	}

	return fmt.Sprintf("%s%s%s%s", dl, m[0], strings.Join(a, ","), dr)
}
