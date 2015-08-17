package urit

import (
	"fmt"
	"net/url"
	"regexp"
	"strings"
)

type Expression struct {
	Match string

	Operator     Operator
	VariableList []string
}

var experreg = regexp.MustCompile(`({([+#.\/;?&,])?([a-z_,]+)})`)

func Split(u URI) []*Expression {
	matches := experreg.FindAllStringSubmatch(string(u), -1)
	if len(matches) == 0 {
		return nil
	}

	var exprs []*Expression
	for _, v := range matches {
		exprs = append(exprs, &Expression{
			Match:        v[1],
			Operator:     Operator(v[2]),
			VariableList: strings.Split(v[3], ","),
		})
	}

	return exprs
}

// Expand expands an expression given the variables. if pact is false the string
// will be returned with the unexpaned variable expressions.
func (e Expression) Expand(pact bool, vars Variables) string {
	var arr []string
	var skp []string // skipped variables

	op := e.Operator

	// switch operator for those that differ in prefix vs. join string
	switch op {
	case "?":
		op = "&"
	case "#":
		op = ","
	}

	for _, v := range e.VariableList {
		val, ok := vars[v]
		if !ok {
			if !pact {
				skp = append(skp, v)
			}

			continue
		}

		// escape through String() to ensure spaces come out as %20 and not +
		val = (&url.URL{Path: val}).String()

		switch op {
		case ";", "?", "&":
			frmt := "%s=%s"
			if val == "" && op == ";" {
				frmt = "%s%s"
			}

			val = fmt.Sprintf(frmt, v, val)
		}

		// insert remaning unexpanded variables as expression
		if len(skp) > 0 {
			arr = append(arr, op.NewExpression(skp))
			skp = nil
		}

		arr = append(arr, val)
	}

	// insert remaning unexpanded variables as expression
	if len(skp) > 0 {
		arr = append(arr, op.NewExpression(skp))
	}

	if len(arr) == 0 {
		return ""
	}

	// final combine is done using the original Operator of the expression
	return e.Operator.Join(arr)
}
