package urit

import (
	"strings"
)

type (
	URI string

	Variables map[string]string
)

func (u URI) Expand(vars ...Variables) URI {
	s := string(u)

	for _, v := range Split(u) {
		s = strings.Replace(s, v.Match, v.Expand(vars...), -1)
	}

	return URI(s)
}
