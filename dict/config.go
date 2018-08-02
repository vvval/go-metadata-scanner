package dict

import (
	"strings"
)

type config struct {
	known    map[string][]string
	booleans []string
	lists    []string
}

func (conf config) Find(name string) (tag, bool) {
	for key, list := range conf.known {
		if strings.EqualFold(key, name) {
			return found(name, key, list)
		}

		for _, val := range list {
			if tagEquals(val, name) {
				return found(name, key, list)
			}
		}
	}

	return notFound(name)
}

func (conf config) IsBoolean(tag tag) bool {
	return oneOf(tag, conf.booleans)
}

func (conf config) IsList(tag tag) bool {
	return oneOf(tag, conf.lists)
}

func found(name, key string, list []string) (tag, bool) {
	t := tag{key, name, list}

	return t, true
}

func notFound(name string) (tag, bool) {
	t := tag{original: name}

	return t, false
}

func oneOf(tag tag, set []string) bool {
	for _, val := range set {
		if tag.has(val) {
			return true
		}
	}

	return false
}

// Name or truncated name (without <GROUP:> prefix) equals
// Tag argument is a full tag name
// Name argument is a searchable input+
func tagEquals(tag, name string) bool {
	return strings.EqualFold(tag, name) || strings.EqualFold(tag, truncatePrefix(name))
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
