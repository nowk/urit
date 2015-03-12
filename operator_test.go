package urit

import (
	"testing"

	"gopkg.in/nowk/assert.v2"
)

func TestNewExpressionBuildsExpressionStringFromArrayStrings(t *testing.T) {
	for _, v := range []struct {
		o, e string
	}{
		{"", "{var,x}"},
		{"+", "{var,x}"},
		{"#", "{#var,x}"},
		{".", "{.var,x}"},
		{"/", "{/var,x}"},
		{";", "{;var,x}"},
		{"?", "{?var,x}"},
		{"&", "{&var,x}"},

		{",", "{,var,x}"},
	} {
		assert.Equal(t, v.e, Operator(v.o).NewExpression([]string{
			"var",
			"x",
		}))
	}
}
