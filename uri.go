package urit

import (
	"strings"
)

type (
	URI string

	Variables map[string]string
)

func (u URI) Expand(vars ...Variables) (URI, error) {
	s := string(u)

	for _, v := range Split(u) {
		str, _ := v.Expand(vars...)
		s = strings.Replace(s, v.Match, str, -1)
	}

	return URI(s), nil
}
