package writecmd

import (
	"fmt"
	"github.com/vvval/go-metadata-scanner/vars/metadata"
	"path/filepath"
	"reflect"
)

type Job struct {
	filename string
	payload  metadata.Payload
}

func (j *Job) Filename() string {
	return j.filename
}

func (j *Job) Payload() metadata.Payload {
	return j.payload
}

func (j *Job) HasPayload() bool {
	return len(j.payload.Tags()) != 0
}

func (j *Job) MergePayload(t metadata.Tags) {
	for _, tag := range j.payload.ListTags() {
		scanned := interface2Slice(t[tag])
		j.payload.UpdateList(tag, scanned)
	}
}

func interface2Slice(slice interface{}) []string {
	value := reflect.ValueOf(slice)
	output := make([]string, value.Len())

	if value.Kind() == reflect.Slice {
		for i := 0; i < value.Len(); i++ {
			output[i] = fmt.Sprintf("%s", value.Index(i))
		}
	}

	return output
}

func NewJob(filename string, payload metadata.Payload) *Job {
	return &Job{filepath.ToSlash(filename), payload}
}
