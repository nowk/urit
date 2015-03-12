package urit

import (
	"fmt"
	"regexp"
	"strings"
)

type Expression struct {
	Match string

	Operator     Operator
	VariableList []string
}

var experreg = regexp.MustCompile(`({([+#.\/;?&,])?([a-z,]+)})`)

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

// Expand expands an expression given the variables
func (e Expression) Expand(d bool, v Variables) string {
	var a []string
	var s []string // skipped variables

	op := e.Operator

	// switch operator for those that differ in prefix vs. join string
	switch op {
	case "?":
		op = "&"
	case "#":
		op = ","
	}

	for _, k := range e.VariableList {
		val, ok := v[k]
		if !ok {
			if !d {
				s = append(s, k)
			}

			continue
		}

		switch op {
		case ";", "?", "&":
			frmt := "%s=%s"
			if val == "" && op == ";" {
				frmt = "%s%s"
			}

			val = fmt.Sprintf(frmt, k, val)
		}

		// insert remaning unexpanded variables as expression
		if len(s) > 0 {
			a = append(a, op.NewExpression(s))
			s = nil
		}

		a = append(a, val)
	}

	// insert remaning unexpanded variables as expression
	if len(s) > 0 {
		a = append(a, op.NewExpression(s))
	}

	if len(a) == 0 {
		return ""
	}

	// final combine is done using the original Operator of the expression
	return e.Operator.Join(a)
}
