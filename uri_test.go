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

// Tests based on http://tools.ietf.org/html/rfc6570

type compare struct {
	a, b URI
}

func TestExpandLevel1SimpleStringExpansion(t *testing.T) {
	for _, v := range []compare{
		{"{var}", "value"},
		{"{hello}", "Hello World!"},
	} {
		assert.Equal(t, v.b, v.a.Expand(vars))
	}
}

func TestExpandLevel2ReservedExpansion(t *testing.T) {
	for _, v := range []compare{
		{"{+var}", "value"},
		{"{+hello}", "Hello World!"},
		{"{+path}/here", "/foo/bar/here"},
		{"here?ref={+path}", "here?ref=/foo/bar"},
	} {
		assert.Equal(t, v.b, v.a.Expand(vars))
	}
}

func TestExpandLevel2FragmentExpansion(t *testing.T) {
	for _, v := range []compare{
		{"X{#var}", "X#value"},
		{"X{#hello}", "X#Hello World!"},
	} {
		assert.Equal(t, v.b, v.a.Expand(vars))
	}
}

func TestExpandLevel3StringExpansion(t *testing.T) {
	for _, v := range []compare{
		{"map?{x,y}", "map?1024,768"},
		{"{x,hello,y}", "1024,Hello World!,768"},
	} {
		assert.Equal(t, v.b, v.a.Expand(vars))
	}
}

func TestExpandLevel3ReservedExpansion(t *testing.T) {
	for _, v := range []compare{
		{"{+x,hello,y}", "1024,Hello World!,768"},
		{"{+path,x}/here", "/foo/bar,1024/here"},
	} {
		assert.Equal(t, v.b, v.a.Expand(vars))
	}
}

func TestExpandLevel3FragmentExpansion(t *testing.T) {
	for _, v := range []compare{
		{"{#x,hello,y}", "#1024,Hello World!,768"},
		{"{#path,x}/here", "#/foo/bar,1024/here"},
	} {
		assert.Equal(t, v.b, v.a.Expand(vars))
	}
}

func TextExpandLevel3LabelExpansion(t *testing.T) {
	for _, v := range []compare{
		{"X{.var}", "X.value"},
		{"X{.x,y}", "X.1024.768"},
		{},
	} {
		assert.Equal(t, v.b, v.a.Expand(vars))
	}
}

func TestExpandLevel3PathSegments(t *testing.T) {
	for _, v := range []compare{
		{"{/var}", "/value"},
		{"{/var,x}/here", "/value/1024/here"},
	} {
		assert.Equal(t, v.b, v.a.Expand(vars))
	}
}

func TestExpandLevel3PathStyleParameters(t *testing.T) {
	for _, v := range []compare{
		{"{;x,y}", ";x=1024;y=768"},
		{"{;x,y,empty}", ";x=1024;y=768;empty"},
	} {
		assert.Equal(t, v.b, v.a.Expand(vars))
	}
}

func TestExpandLevel3FormStyleQuery(t *testing.T) {
	for _, v := range []compare{
		{"{?x,y}", "?x=1024&y=768"},
		{"{?x,y,empty}", "?x=1024&y=768&empty="},
	} {
		assert.Equal(t, v.b, v.a.Expand(vars))
	}
}

func TestExpandLevel3FormStyleQueryContinuation(t *testing.T) {
	for _, v := range []compare{
		{"?fixed=yes{&x}", "?fixed=yes&x=1024"},
		{"{&x,y,empty}", "&x=1024&y=768&empty="},
	} {
		assert.Equal(t, v.b, v.a.Expand(vars))
	}
}

func TestExpandDoesNotExpandIfNoVariableValue(t *testing.T) {
	u := URI("{/var,x}/here{?x,y}").Expand(Variables{
		"var": "/foo/bar",
		"y":   "768",
	})
	assert.Equal(t, URI("/foo/bar/here?y=768"), u)
}

func TestInspectAllowsForContinuousReexpanding(t *testing.T) {
	u := URI("{/var,x}/here{?x,y}").Inspect(Variables{
		"var": "/foo/bar",
	})
	assert.Equal(t, URI("/foo/bar{/x}/here?{&x,y}"), u)

	u = u.Inspect(Variables{
		"x": "1024",
		"y": "768",
	})
	assert.Equal(t, URI("/foo/bar/1024/here?&x=1024&y=768"), u)
	assert.Equal(t, URI("/foo/bar/1024/here?&x=1024&y=768"), u.Expand(nil))
}

func TestFixMultipleInspectCausesDoubleAmperOnQueryExpand(t *testing.T) {
	var u URI = "/users{/username}/repos{?type,sort,direction}"

	u = u.Inspect(Variables{
		"username": "nowk",
	})

	u = u.Inspect(Variables{
		"sort": "created",
	})

	assert.Equal(t, URI("/users/nowk/repos?&sort=created"), u.Expand(nil))
}
