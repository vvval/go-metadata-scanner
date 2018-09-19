package writecmd

import (
	"fmt"
	"github.com/vvval/go-metadata-scanner/vars/metadata"
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
	return j.payload.Tags().Count() != 0
}

func (j *Job) MergePayload(tags metadata.Tags) {
	for _, tag := range j.payload.ListTags() {
		if val, ok := tags.Tag(tag); ok {
			scanned := interface2Slice(val)
			j.payload.UpdateList(tag, scanned)
		}
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
	return &Job{filename, payload}
}
