package urit

import (
	"testing"

	"gopkg.in/nowk/assert.v2"
)

var vars = Variables{
	"var":   "value",
	"hello": "Hello World!",
	"empty": "",
	"path":  "/foo/bar",
	"x":     "1024",
	"y":     "768",
}

func TestSplit(t *testing.T) {
	a := Split("{+x,hello,y}")
	assert.Equal(t, []*Expression{
		{
			"{+x,hello,y}",
			"+",
			"x,hello,y",
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

// Tests based on http://tools.ietf.org/html/rfc6570

type compare struct {
	a, b URI
}

func TestExpandLevel1SimpleStringExpansion(t *testing.T) {
	for _, v := range []compare{
		{"{var}", "value"},
		{"{hello}", "Hello World!"},
	} {
		exp, err := v.a.Expand(vars)
		if err != nil {
			t.Fatal(err)
		}
		assert.Equal(t, v.b, exp)
	}
}

func TestExpandLevel2ReservedExpansion(t *testing.T) {
	for _, v := range []compare{
		{"{+var}", "value"},
		{"{+hello}", "Hello World!"},
		{"{+path}/here", "/foo/bar/here"},
		{"here?ref={+path}", "here?ref=/foo/bar"},
	} {
		exp, err := v.a.Expand(vars)
		if err != nil {
			t.Fatal(err)
		}
		assert.Equal(t, v.b, exp)
	}
}

func TestExpandLevel2FragmentExpansion(t *testing.T) {
	for _, v := range []compare{
		{"X{#var}", "X#value"},
		{"X{#hello}", "X#Hello World!"},
	} {
		exp, err := v.a.Expand(vars)
		if err != nil {
			t.Fatal(err)
		}
		assert.Equal(t, v.b, exp)
	}
}

func TestExpandLevel3StringExpansion(t *testing.T) {
	for _, v := range []compare{
		{"map?{x,y}", "map?1024,768"},
		{"{x,hello,y}", "1024,Hello World!,768"},
	} {
		exp, err := v.a.Expand(vars)
		if err != nil {
			t.Fatal(err)
		}
		assert.Equal(t, v.b, exp)
	}
}

func TestExpandLevel3ReservedExpansion(t *testing.T) {
	for _, v := range []compare{
		{"{+x,hello,y}", "1024,Hello World!,768"},
		{"{+path,x}/here", "/foo/bar,1024/here"},
	} {
		exp, err := v.a.Expand(vars)
		if err != nil {
			t.Fatal(err)
		}
		assert.Equal(t, v.b, exp)
	}
}

func TestExpandLevel3FragmentExpansion(t *testing.T) {
	for _, v := range []compare{
		{"{#x,hello,y}", "#1024,Hello World!,768"},
		{"{#path,x}/here", "#/foo/bar,1024/here"},
	} {
		exp, err := v.a.Expand(vars)
		if err != nil {
			t.Fatal(err)
		}
		assert.Equal(t, v.b, exp)
	}
}

func TextExpandLevel3LabelExpansion(t *testing.T) {
	for _, v := range []compare{
		{"X{.var}", "X.value"},
		{"X{.x,y}", "X.1024.768"},
		{},
	} {
		exp, err := v.a.Expand(vars)
		if err != nil {
			t.Fatal(err)
		}
		assert.Equal(t, v.b, exp)
	}
}

func TestExpandLevel3PathSegments(t *testing.T) {
	for _, v := range []compare{
		{"{/var}", "/value"},
		{"{/var,x}/here", "/value/1024/here"},
	} {
		exp, err := v.a.Expand(vars)
		if err != nil {
			t.Fatal(err)
		}
		assert.Equal(t, v.b, exp)
	}
}

func TestExpandLevel3PathStyleParameters(t *testing.T) {
	for _, v := range []compare{
		{"{;x,y}", ";x=1024;y=768"},
		{"{;x,y,empty}", ";x=1024;y=768;empty"},
	} {
		exp, err := v.a.Expand(vars)
		if err != nil {
			t.Fatal(err)
		}
		assert.Equal(t, v.b, exp)
	}
}

func TestExpandLevel3FormStyleQuery(t *testing.T) {
	for _, v := range []compare{
		{"{?x,y}", "?x=1024&y=768"},
		{"{?x,y,empty}", "?x=1024&y=768&empty="},
	} {
		exp, err := v.a.Expand(vars)
		if err != nil {
			t.Fatal(err)
		}
		assert.Equal(t, v.b, exp)
	}
}

func TestExpandLevel3FormStyleQueryContinuation(t *testing.T) {
	for _, v := range []compare{
		{"?fixed=yes{&x}", "?fixed=yes&x=1024"},
		{"{&x,y,empty}", "&x=1024&y=768&empty="},
	} {
		exp, err := v.a.Expand(vars)
		if err != nil {
			t.Fatal(err)
		}
		assert.Equal(t, v.b, exp)
	}
}

func TextExpandMultipleExpressions(t *testing.T) {
	for _, v := range []compare{
		{"/a{/var,x}/here?fixed=yes{&x}{&y}", "/a/value/1024/here?fixed=yes&x=1024&y=768"},
	} {
		exp, err := v.a.Expand(vars)
		if err != nil {
			t.Fatal(err)
		}
		assert.Equal(t, v.b, exp)
	}
}
