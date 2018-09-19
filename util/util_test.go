package util

import (
	"testing"
)

func TestPathsEqual(t *testing.T) {
	type check struct {
		a, b string
		exp  bool
	}

	set := []check{
		{"file\\with/different/slashes\\usage", "file/with\\different/slashes/usage", true},
		{"file\\with/different/slashes\\usage", "file/with/different/slashes/usage", true},
		{"file/with/different/slashes/usage", "file\\with\\different\\slashes\\usage", true},
		{"file\\with\\different\\slashes\\usage", "file\\with\\different\\slashes\\usage", true},
		{"file/with/different/slashes/usage", "file/with/different/slashes/usage", true},
		{"file/with/different/slashes/usage", "file/with/different/slashes/usage/", true},
		{"file/with/different/slashes/usage\\", "file/with/different/slashes/usage/", true},
		{"file/with/different/slashes/usage\\", "file/with/different/slashes/usage", true},
		{"file/with/different/slashes/usage\\", "file/with/different/slashes/usage/oops", false},
	}

	for i, s := range set {
		c := PathsEqual(s.a, s.b)
		if c != s.exp {
			t.Errorf("values compare failed for %d:\ninput `%s`, `%s`\ngot `%t` \nexpected `%t`", i, s.a, s.b, c, s.exp)
		}
	}
}

func TestTokenizer(t *testing.T) {
	set := map[string][]string{
		"a,b,,c, d": {"a", "b", "c", "d"},
		"a,b,a":     {"a", "b", "a"},
		`a,"b,c",d`: {"a", "d", `b,c`},
		"a;b,c":     {"a", "b", "c"},
	}

	for str, exp := range set {
		tokens := SplitKeywords(str)
		if !Equal(exp, tokens) {
			t.Errorf("tokens not equal:\ngot `%s`\nexpected `%s`", tokens, exp)
		}
	}
}

func TestAdjustSize(t *testing.T) {
	type check struct {
		n, d, min, an, ad int
	}
	set := []check{
		{10, 3, 5, 2, 5},
		{10, 6, 5, 2, 5},
		{10, 3, 3, 3, 4},
		{10, 3, 2, 3, 4},
		{10, 3, 11, 1, 11},
		{10, 3, 9, 2, 9},
	}

	for i, v := range set {
		p, c := AdjustSizes(v.n, v.d, v.min)
		if p != v.an || c != v.ad {
			t.Errorf("values are not equal for %d:\ninput `%d`, `%d` and `%d`\ngot `%d` and `%d`\nexpected `%d` and `%d`", i, v.n, v.d, v.min, p, c, v.an, v.ad)
		}
	}
}

func TestExtension(t *testing.T) {
	set := [][]string{
		{"filename.ext", "ext"},
		{".ext", "ext"},
		{"filename", ""},
		{"filename.", ""},
	}

	for i, str := range set {
		ext := Extension(str[0])
		if ext != str[1] {
			t.Errorf("extensions not equal for %d:\ngot `%s`\nexpected `%s`", i, ext, str[1])
		}
	}
}

func TestEqual(t *testing.T) {
	type check struct {
		a, b []string
		exp  bool
	}

	set := []check{
		{[]string{"a", "b"}, []string{"b", "a"}, true},
		{[]string{"a", "b"}, []string{"b", "c"}, false},
		{[]string{"a", "b"}, []string{"b", "A"}, false},
		{[]string{"a", "b"}, []string{"b"}, false},
	}

	for i, v := range set {
		res := Equal(v.a, v.b)
		if res != v.exp {
			t.Errorf("equality failed for %d:\ngot `%t` \nexpected `%t`", i, res, v.exp)
		}
	}
}

func TestUnique(t *testing.T) {
	type check struct {
		a, b []string
	}

	set := []check{
		{[]string{"a", "b", "c"}, []string{"a", "b", "c"}},
		{[]string{"a", "b", "a"}, []string{"b", "a"}},
		{[]string{"a", "B", "b"}, []string{"b", "a", "B"}},
	}

	for i, v := range set {
		res := UniqueValues(v.a)
		if !Equal(res, v.b) {
			t.Errorf("equality failed for %d:\ngot `%s` \nexpected `%s`", i, res, v.b)
		}
	}
}
