package util

import (
	"reflect"
	"sort"
	"strings"
)

func UniqueValues(value []string) []string {
	var m = make(map[string]bool)
	var out []string
	for _, v := range value {
		v = strings.Trim(v, " ")
		if _, ok := m[v]; !ok && len(v) > 0 {
			m[v] = true
			out = append(out, v)
		}
	}

	return out
}

func Equal(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}

	sort.Strings(a)
	sort.Strings(b)

	return reflect.DeepEqual(a, b)
}
