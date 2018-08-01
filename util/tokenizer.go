package util

import (
	"strings"
	"unicode"
)

var (
	separators = &unicode.RangeTable{
		R16: []unicode.Range16{
			{0x002c, 0x002c, 1}, //comma
			{0x003b, 0x003b, 1}, //semicolon
		},
	}
	quotationMarks = &unicode.RangeTable{
		R16: []unicode.Range16{
			{0x0022, 0x0022, 1}, //quot
		},
	}
	trimRunes = `,; "`
)

func SplitKeywords(input string) []string {
	var (
		keyword   string
		keywords  []string
		lastIndex int
		quotFound bool
	)

	for index, r := range []rune(input) {
		if unicode.In(r, separators) && !quotFound {
			keyword = fetchKeyword(input, lastIndex, index)
			lastIndex = index

			if len(keyword) != 0 {
				keywords = append(keywords, keyword)
			}

			continue
		}

		if index == len(input)-1 {
			keyword = fetchKeyword(input, lastIndex, len(input))
			lastIndex = index

			if len(keyword) != 0 {
				keywords = append(keywords, keyword)
			}

			continue
		}

		if unicode.In(r, quotationMarks) {
			quotFound = !quotFound
			continue
		}
	}

	return keywords
}

func fetchKeyword(input string, start, end int) string {
	return trim(cut(input, start, end))
}

func trim(input string) string {
	return strings.Trim(input, trimRunes)
}

func cut(input string, start, end int) string {
	return string([]rune(input)[start:end])
}
