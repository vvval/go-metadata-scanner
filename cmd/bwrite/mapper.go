package bwrite

import (
	"github.com/vvval/go-metadata-scanner/cmd/config"
	"github.com/vvval/go-metadata-scanner/cmd/metadata"
	"github.com/vvval/go-metadata-scanner/dict"
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
	d := dict.Get()

	for index, value := range input {
		key, ok := columns[index]
		if !ok {
			// Unmapped key, skip
			continue
		}

		t, found := d.Find(key)
		if !found {
			// Unknown tag, skip
			continue
		}

		for _, tag := range t.Map() {
			if d.IsBoolean(t) {
				line.AddTag(tag, len(value) != 0)
			} else if len(value) != 0 {
				if d.IsList(t) {
					line.AddTag(tag, util.SplitKeywords(value))
				} else {
					line.AddTag(tag, value)
				}
			}
		}
	}

	return line
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
