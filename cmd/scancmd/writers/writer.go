package writers

import (
	"fmt"
	"github.com/vvval/go-metadata-scanner/vars"
	"os"
	"strings"
)

type Writer interface {
	Open(filename string, headers []string) error
	Close() error
	Write(file *vars.File) error
}

type BaseWriter struct {
	file     *os.File
	filename string
	headers  []string
}

func NewWriter(filename string, headers []string) BaseWriter {
	return BaseWriter{
		filename: filename,
		headers:  headers,
	}
}

func openFile(filename string) (*os.File, error) {
	file, err := os.Create(filename)
	if err != nil {
		return nil, err
	}

	return file, nil
}

func closeFile(file *os.File) {
	if file != nil {
		file.Close()
	}
}

func tagsByGroups(f *vars.File, groups []string) map[string]map[string]string {
	tags := make(map[string]map[string]string)
	for tag, value := range f.Tags() {
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
