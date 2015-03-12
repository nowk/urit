package urit

import (
	"testing"

	"gopkg.in/nowk/assert.v2"
)

func TestSplitReturnsExpressions(t *testing.T) {
	a := Split("{+x,hello,y}")
	assert.Equal(t, []*Expression{
		{
			"{+x,hello,y}",
			"+",
			[]string{
				"x",
				"hello",
				"y",
			},
		},
	}, a)
}

func TestSplitReturnsNilOnNonMatch(t *testing.T) {
	for _, v := range []string{
		"",
		"+x,hello,y",
		"{}",
		"{+x,hello,y",
		"+x,hello,y}",
	} {
		a := Split(URI(v))
		assert.Nil(t, a)
	}
}

func TestInspectionUnexpandedPath(t *testing.T) {
	e := Expression{
		Operator: "/",
		VariableList: []string{
			"var",
			"x",
			"y",
			"z",
		},
	}

	for _, v := range []struct {
		s string
		v Variables
	}{
		{"{/var}/1024{/y}/1", Variables{"x": "1024", "z": "1"}},
		{"/value{/x}/768{/z}", Variables{"var": "value", "y": "768"}},
		{"{/var}/1024/768{/z}", Variables{"x": "1024", "y": "768"}},
		{"/value{/x,y}/1", Variables{"var": "value", "z": "1"}},
		{"{/var,x}/768{/z}", Variables{"y": "768"}},
		{"{/var}/1024{/y,z}", Variables{"x": "1024"}},
	} {
		assert.Equal(t, v.s, e.Expand(false, v.v))
	}
}

func TestInspectionUnexpandedQuery(t *testing.T) {
	e := Expression{
		Operator: "?",
		VariableList: []string{
			"var",
			"x",
			"y",
			"z",
			"a",
		},
	}

	for _, v := range []struct {
		s string
		v Variables
	}{
		{"?var=value{&x}&y=768{&z,a}", Variables{"var": "value", "y": "768"}},
		{"?{&var,x}&y=768{&z}&a=2", Variables{"y": "768", "a": "2"}},
	} {
		assert.Equal(t, v.s, e.Expand(false, v.v))
	}
}

func TestInspectionUnexpandedHash(t *testing.T) {
	e := Expression{
		Operator: "#",
		VariableList: []string{
			"var",
			"x",
			"y",
			"z",
			"a",
		},
	}

	for _, v := range []struct {
		s string
		v Variables
	}{
		{"#value{,x},768{,z,a}", Variables{"var": "value", "y": "768"}},
		{"#{,var,x},768{,z},2", Variables{"y": "768", "a": "2"}},
	} {
		assert.Equal(t, v.s, e.Expand(false, v.v))
	}
}
