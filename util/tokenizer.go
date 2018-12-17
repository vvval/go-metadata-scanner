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
	trimRunes = ",; \"\t\n"
)

type tokenizer struct {
	keyword   string
	keywords  []string
	lastIndex int
	quotFound bool
}

var t tokenizer

func (t *tokenizer) addKeyword(input []rune, index int) {
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

	runes := []rune(input)
	for index, r := range runes {
		if unicode.In(r, separators) && !t.quotFound {
			t.addKeyword(runes, index)

			continue
		}

		//end of string
		if index == len(runes)-1 {
			t.addKeyword(runes, len(runes))

			continue
		}

		if unicode.In(r, quotationMarks) {
			t.toggleQuot()

			continue
		}
	}

	return t.keywords
}

func fetchKeyword(runes []rune, start, end int) string {
	return trim(cut(runes, start, end))
}

func trim(s string) string {
	return strings.Trim(s, trimRunes)
}

func cut(runes []rune, start, end int) string {
	if end > len(runes) {
		end = len(runes)
	}

	return string(runes[start:end])
}
