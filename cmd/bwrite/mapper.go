package bwrite

import (
	"github.com/vvval/go-metadata-scanner/cmd/config"
	"regexp"
	"strings"
)

const separator string = "<sep>"

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
func MapLineToColumns(columns map[int]string, line []string) map[string]string {
	output := map[string]string{}

	for index, value := range line {
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
			if IsListTag(tag) {
				value = filterEmptyStringChunks(value)
			}

			output[tag] = value
		}
	}

	return output
}

func IsListTag(name string) bool {
	name = findTagAliasKey(name)
	for _, tag := range config.AppConfig().ListTags {
		if strings.EqualFold(name, tag) || strings.EqualFold(name, truncateKeyPrefix(tag)) {
			// Given name is a presented alias, or name (ignoring <group:> prefix)
			return true
		}
	}

	return false
}

func filterEmptyStringChunks(s string) string {
	var regex = regexp.MustCompile(`\s?[,;]\s?`)
	s = regex.ReplaceAllString(s, separator)

	filtered := trimEmptyChunks(strings.Split(s, separator))

	return strings.Join(filtered, separator)
}

func trimEmptyChunks(value []string) []string {
	var chunks []string

	for _, chunk := range value {
		chunk = strings.Trim(chunk, " ")
		if len(chunk) == 0 {
			continue
		}

		chunks = append(chunks, chunk)
	}

	return chunks
}

// Map columns to a known tag map
// Skip 1st column (dedicated to a file names) and empty columns
func MapColumns(columns []string) map[int]string {
	output := map[int]string{}
	for index, value := range columns {
		value = strings.Trim(value, " ")
		if index == 0 || len(value) == 0 {
			// Skip 1st column and empty columns
			continue
		}

		output[index] = findTagAliasKey(value)
	}

	return output
}

func findTagAliasKey(name string) string {
	for key, values := range config.AppConfig().TagMap {
		if strings.EqualFold(name, key) {
			// Tag matched the map key
			return key
		}

		for _, value := range values {
			if strings.EqualFold(name, value) || strings.EqualFold(name, truncateKeyPrefix(value)) {
				// Tag matched one of map values
				return key
			}
		}
	}

	return name
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
