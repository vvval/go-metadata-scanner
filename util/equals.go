package util

import "strings"

// Name or truncated name (without <GROUP:> prefix) equals
// Tag argument is a full tag name
// Name argument is a searchable input+
func TagEquals(tag, name string) bool {
	return Equals(tag, name) || Equals(tag, truncatePrefix(name))
}

// Short alias
func Equals(s1, s2 string) bool {
	return strings.EqualFold(s1, s2)
}

// Cut <group:> prefix if found
func truncatePrefix(tag string) string {
	prefixEnding := strings.Index(tag, ":")
	if prefixEnding == -1 {
		return tag
	}

	runes := []rune(tag)

	return string(runes[prefixEnding+1:])
}
