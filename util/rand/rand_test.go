package rand

import (
	"regexp"
	"testing"
)

func TestRand(t *testing.T) {
	set := []int{10, 100}
	var reg = &regexp.Regexp{}
	reg = regexp.MustCompile("^[a-zA-Z0-9]*$")
	for _, n := range set {
		str := Strings(n)
		if len(str) != n || !reg.MatchString(str) {
			t.Errorf("random string incorrect:\ngot `%s` of length `%d`\nexpected a-zA-A0-9 regex of `%d` length", str, len(str), n)
		}
	}
}
