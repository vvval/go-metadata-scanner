package dict

import "github.com/vvval/go-metadata-scanner/util"

type tag struct {
	key, original string
	tags          []string
}

func (t tag) Map() []string {
	return t.tags
}

// Check if val is a key or one of key mapped values
func (t tag) has(name string) bool {
	if util.Equals(name, t.key) {
		return true
	}

	if t.inMap(name) {
		return true
	}

	return false
}

func (t tag) inMap(name string) bool {
	for _, val := range t.tags {
		if util.Equals(name, val) {
			return true
		}
	}

	return false
}
