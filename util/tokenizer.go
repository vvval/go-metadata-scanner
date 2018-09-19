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

type tokenizer struct {
	keyword   string
	keywords  []string
	lastIndex int
	quotFound bool
}

var t tokenizer

func (t *tokenizer) addKeyword(input string, index int) {
	keyword := fetchKeyword(input, t.lastIndex, index)
	t.lastIndex = index

	if len(keyword) != 0 {
		t.append(keyword)
	}
}

func (t *tokenizer) append(keyword string) {
	t.keywords = append(t.keywords, keyword)
}

func (t *tokenizer) toggleQuot() {
	t.quotFound = !t.quotFound
}

func SplitKeywords(input string) []string {
	t = tokenizer{}

	for index, r := range []rune(input) {
		if unicode.In(r, separators) && !t.quotFound {
			t.addKeyword(input, index)

			continue
		}

		if index == len(input)-1 {
			t.addKeyword(input, len(input))

			continue
		}

		if unicode.In(r, quotationMarks) {
			t.toggleQuot()

			continue
		}
	}

	return t.keywords
}

func fetchKeyword(input string, start, end int) string {
	return trim(cut(input, start, end))
}

func trim(input string) string {
	return strings.Trim(input, trimRunes)
}

func cut(input string, start, end int) string {
	r := []rune(input)
	if end > len(r) {
		end = len(r)
	}

	return string(r[start:end])
}
