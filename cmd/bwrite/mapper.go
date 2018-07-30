package bwrite

import (
	"github.com/vvval/go-metadata-scanner/cmd/config"
	"strings"
)

// Map columns to a known tag map
func MapColumns(line []string) map[int]string {
	output := map[int]string{}
	for key, values := range config.AppConfig().TagMap {
		for index, name := range line {
			name = strings.Trim(name, " ")
			// Skip empty lines and 1st column
			if index == 0 || len(name) == 0 {
				continue
			}

			// Tag map key matches
			if strings.EqualFold(name, key) {
				output[index] = key

				continue
			}

			// Tag map value matches
			for _, value := range values {
				if strings.EqualFold(name, value) || strings.EqualFold(name, truncateKeyPrefix(value)) {
					output[index] = key

					break
				}
			}
		}
	}

	return output
}

func MapLine(columns map[int]string, data []string) map[string]string {
	output := map[string]string{}

	for index, value := range data {
		key, ok := columns[index]
		// Unmapped key, skip
		if !ok {
			continue
		}

		// Unknown tag
		tags, ok := (config.AppConfig().TagMap)[key]
		if !ok {
			continue
		}

		for _, tag := range tags {
			output[tag] = value
		}
	}

	return output
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
