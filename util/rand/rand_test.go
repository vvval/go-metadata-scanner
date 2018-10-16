package rand

import (
	"regexp"
	"testing"
	"time"
)

func TestRand(t *testing.T) {
	set := []int{10, 100}
	var reg = &regexp.Regexp{}
	reg = regexp.MustCompile("^[a-zA-Z0-9]*$")
	for i, n := range set {
		str := Strings(n)
		if len(str) != n || !reg.MatchString(str) {
			t.Errorf("random string incorrect (line `%d`):\ngot `%s` of length `%d`\nexpected `a-zA-A0-9` regex of `%d` length", i, str, len(str), n)
		}
	}

	str1 := Strings(10)
	time.Sleep(time.Microsecond)
	str2 := Strings(10)

	if str1 == str2 {
		t.Errorf("random strings should not repeat:\ngot `%s` `%s`", str1, str2)
	}
}
