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
	if tag, found := inKnown(name); found {
		return tag, found
	}

	for _, b := range conf.booleans {
		if strings.EqualFold(b, name) {
			return found("", name, []string{name})
		}
	}

	for _, l := range conf.lists {
		if strings.EqualFold(l, name) {
			return found("", name, []string{name})
		}
	}

	return notFound(name)
}

func inKnown(name string) (tag, bool) {
	for key, list := range conf.known {
		if strings.EqualFold(key, name) {
			return found(key, name, list)
		}

		for _, val := range list {
			if tagEquals(val, name) {
				return found(key, name, list)
			}
		}
	}

	return notFound(name)
}

func found(key, name string, list []string) (tag, bool) {
	t := tag{key, name, list}

	return t, true
}

func notFound(name string) (tag, bool) {
	t := tag{original: name}

	return t, false
}

func (conf config) IsBoolean(key, tag string) bool {
	return oneOf(key, tag, conf.booleans)
}

func (conf config) IsList(key, tag string) bool {
	return oneOf(key, tag, conf.lists)
}

func oneOf(key, tag string, set []string) bool {
	for _, val := range set {
		if strings.EqualFold(tag, val) || strings.EqualFold(key, val) {
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
