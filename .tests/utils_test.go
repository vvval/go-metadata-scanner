package _tests

import (
	"github.com/vvval/go-metadata-scanner/util"
	"reflect"
	"testing"
)

func TestTokenizer(t *testing.T) {
	set := map[string][]string{
		"a,b,,c, d": {"a", "b", "c", "d"},
		"a,b,a":     {"a", "b", "a"},
		`a,"b,c",d`: {"a", `b,c`, "d"},
		"a;b,c":     {"a", "b", "c"},
	}

	for str, exp := range set {
		tokens := util.SplitKeywords(str)
		if !reflect.DeepEqual(exp, tokens) {
			t.Errorf("tokens not equal:\ngot `%s`\nexpected `%s`", tokens, exp)
		}
	}
}

func TestAdjustSize(t *testing.T) {
	set := [][]int{
		{10, 3, 5, 2, 5},
		{10, 6, 5, 2, 5},
		{10, 3, 3, 3, 4},
		{10, 3, 2, 3, 4},
		{10, 3, 11, 1, 11},
		{10, 3, 9, 2, 9},
	}

	for _, v := range set {
		p, c := util.AdjustSizes(v[0], v[1], v[2])
		if p != v[3] || c != v[4] {
			t.Errorf("values are not equal:\ninput `%d`, `%d` and `%d`\ngot `%d` and `%d`\nexpected `%d` and `%d`", v[0], v[1], v[2], p, c, v[3], v[4])
		}
	}
}
