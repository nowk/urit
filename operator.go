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
func (o Operator) Join(arr []string) string {
	var str string

	switch o {
	case "/":
		for i, v := range arr {
			if v[:1] != "{" {
				arr[i] = "/" + v
			}
		}
		str = filepath.Join(strings.Join(arr, ""))

	case "#", "?", ".", ";", "&":
		m := opmap[o]
		for i, v := range arr {
			if i > 0 && v[:1] != "{" {
				arr[i] = m[1] + v
			}
		}
		str = m[0] + strings.Join(arr, "")

		// TODO revisit, not particularly pretty
		if str[:2] == "&{" {
			str = str[1:] // remove extra &
		}

	default: // handles +, `{+var}` and blank, `{var}`
		str = strings.Join(arr, ",")
	}

	return str
}

// left and right expression delimiters { }
const (
	delimL = "{"
	delimR = "}"
)

func (o Operator) NewExpression(arr []string) string {
	if len(arr) == 0 {
		return ""
	}

	m, ok := opmap[o]
	if !ok {
		m = []string{"", ","}
	}

	return fmt.Sprintf("%s%s%s%s", delimL, m[0], strings.Join(arr, ","), delimR)
}
