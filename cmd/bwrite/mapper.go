package bwrite

import (
	"github.com/vvval/go-metadata-scanner/cmd/config"
	"github.com/vvval/go-metadata-scanner/cmd/metadata"
	"github.com/vvval/go-metadata-scanner/util"
	"strings"
)

// Script will pick all values for line data that are presented in columns (by column index)
// and fill output array with line value for each tag aliases from TagMap.
//
// Example:
// 		Input line is [1:"keywords values", 3:"some description"
// 		Input columns are [1:keywords,2:title,3:description]
// 		Output is [
// 			"IPTC:keywords1":"keywords values",
// 			"XMP:keywords2":"keywords values",
// 			"IPTC:description1":"some description",
// 			"XMP:description2":"some description"
// 		]
func MapLineToColumns(columns map[int]string, input []string) metadata.Line {
	line := metadata.NewLine()

	for index, value := range input {
		key, ok := columns[index]
		if !ok {
			// Unmapped key, skip
			continue
		}

		tags, ok := (config.AppConfig().TagMap)[key]
		if !ok {
			// Unknown tag, skip
			continue
		}

		for _, tag := range tags {
			var val interface{}

			switch {
			case isListTag(tag):
				val = util.SplitKeywords(value)
			case isBoolTag(tag):
				val = len(value) != 0
			default:
				val = value
			}

			line.AddTag(tag, val)
		}
	}

	return line
}

func isListTag(name string) bool {
	return oneOf(name, config.AppConfig().ListTags)
}

func isBoolTag(name string) bool {
	return oneOf(name, config.AppConfig().BoolTags)
}

func oneOf(name string, set []string) bool {
	key := tagAliasKey(name)
	for _, tag := range set {
		if tagEquals(name, tag) || tagEquals(key, tag) {
			// Given name is a presented alias, or name (ignoring <group:> prefix)
			return true
		}
	}

	return false
}

// Map columns to a known tag map
// Skip 1st column (dedicated to a file names) and empty columns
func MapColumns(columns []string) map[int]string {
	output := map[int]string{}
	for i, tag := range columns {
		tag = strings.Trim(tag, " ")
		if i == 0 || len(tag) == 0 {
			// Skip 1st column and empty columns
			continue
		}

		output[i] = tagAliasKey(tag)
	}

	return output
}

// Find key in tagMap (aliases)
// If not found - return input value
func tagAliasKey(name string) string {
	for key, values := range config.AppConfig().TagMap {
		if equals(name, key) {
			// Tag matched the map key
			return key
		}

		for _, value := range values {
			if tagEquals(name, value) {
				// Tag matched one of map values
				return key
			}
		}
	}

	return name
}

// Key or truncated key equals
func tagEquals(s, t string) bool {
	return equals(s, t) || equals(s, truncateKeyPrefix(t))
}

// Short alias
func equals(s, t string) bool {
	return strings.EqualFold(s, t)
}

// Cut <group:> prefix if found
func truncateKeyPrefix(key string) string {
	prefixEnding := strings.Index(key, ":")
	if prefixEnding == -1 {
		return key
	}

	runes := []rune(key)

	return string(runes[prefixEnding+1:])
}

//func test(s string) {
//	fmt.Printf("-----input `%s`\n", s)
//
//	var (
//		separators = &unicode.RangeTable{
//			R16: []unicode.Range16{
//				{0x002c, 0x002c, 1},
//				{0x003b, 0x003b, 1},
//			},
//		}
//		quotationMark = &unicode.RangeTable{
//			R16: []unicode.Range16{
//				{0x0022, 0x0022, 1},
//			},
//		}
//		keywords         []string
//		lastKeywordIndex int
//		quotFound        bool
//	)
//
//	for i, r := range []rune(s) {
//		if unicode.In(r, separators) {
//			if !quotFound {
//				fmt.Printf("+++++separator pos %d, s: %s\n", i, s)
//				keywords = append(keywords, strings.Trim(string([]rune(s)[lastKeywordIndex:i]), `,; "`))
//				lastKeywordIndex = i
//			}
//
//			continue
//		}
//
//		if i == len(s)-1 {
//			fmt.Printf("+++++separator pos %d, s: %s\n", i, s)
//			keywords = append(keywords, strings.Trim(string([]rune(s)[lastKeywordIndex:len(s)]), `,; "`))
//			lastKeywordIndex = i
//
//			continue
//		}
//
//		if unicode.In(r, quotationMark) {
//			quotFound = !quotFound
//			continue
//		}
//		//fmt.Printf("*****debug: %d, %d\n", i, len(s))
//	}
//	fmt.Printf("=====output: %+v\n")
//	for _, k := range keywords {
//		fmt.Printf("=%+v\n", k)
//	}
//}
//
//func trimEmptyChunks(value []string) []string {
//	var chunks []string
//
//	for _, chunk := range value {
//		chunk = strings.Trim(chunk, " ")
//		if len(chunk) == 0 {
//			continue
//		}
//
//		chunks = append(chunks, chunk)
//	}
//
//	return chunks
//}
