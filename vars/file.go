package vars

import (
	"encoding/json"
	"fmt"
	"github.com/vvval/go-metadata-scanner/vars/metadata"
	"path/filepath"
	"strings"
)

type File struct {
	filename string //parsed with filepath.ToSlashes
	tags     metadata.Tags
}

func NewFile(filename string, tags metadata.Tags) File {
	return File{filepath.ToSlash(filename), tags}
}

func (f *File) Filename() string {
	return f.filename
}

func (f *File) RelPath(base string) string {
	rel, err := filepath.Rel(base, f.filename)
	if err != nil {
		return f.filename
	}

	return rel
}

func (f *File) Tags() metadata.Tags {
	return f.tags
}

func (f *File) Pack(headers []string) map[string]string {
	output := make(map[string]string)

	for header, value := range f.splitTagsToGroups(headers) {
		packed, err := json.Marshal(value)
		if err == nil {
			output[header] = string(packed)
		}
	}

	return output
}

func (f *File) splitTagsToGroups(groups []string) map[string]map[string]string {
	tags := make(map[string]map[string]string)
	for tag, value := range f.tags {
		for _, group := range groups {
			if groupPrefixMatch(tag, group) {
				tagGroup := tagGroup(tag)
				if _, ok := tags[tagGroup]; !ok {
					tags[tagGroup] = make(map[string]string)
				}

				tags[tagGroup][tagName(tag)] = fmt.Sprintf("%v", value)
			}
		}
	}

	return tags
}

func groupPrefixMatch(tag, group string) bool {
	return strings.HasPrefix(strings.ToLower(tag), strings.ToLower(group)+":")
}

func tagGroup(tag string) string {
	return tag[:groupSeparatorPos(tag)]
}
func tagName(tag string) string {
	return tag[groupSeparatorPos(tag)+1:]
}

func groupSeparatorPos(tag string) int {
	return strings.Index(tag, ":")
}
