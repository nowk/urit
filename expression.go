package urit

import (
	"fmt"
	"regexp"
	"strings"
)

var experreg = regexp.MustCompile(`({([+#.\/;?&])?([a-z,]+)})`)

type Expression struct {
	Match string

	Operator     Operator
	VariableList []string
}

func (e Expression) Expand(vars ...Variables) (string, error) {
	var a []string

	for _, k := range e.VariableList {
		for _, v := range vars {
			val, ok := v[k]
			if !ok {
				continue
			}

			switch e.Operator {
			case ";", "?", "&":
				frmt := "%s=%s"
				if val == "" && e.Operator == ";" {
					frmt = "%s%s"
				}

				val = fmt.Sprintf(frmt, k, val)
			}

			a = append(a, val)
		}
	}

	return e.Operator.Join(a)
}

func Split(u URI) []*Expression {
	m := experreg.FindAllStringSubmatch(string(u), -1)
	if len(m) == 0 {
		return nil
	}

	var e []*Expression
	for _, v := range m {
		e = append(e, &Expression{
			Match:        v[1],
			Operator:     Operator(v[2]),
			VariableList: strings.Split(v[3], ","),
		})
	}

	return e
}
