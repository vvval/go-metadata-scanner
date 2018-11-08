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
	if value.Kind() == reflect.Slice {
		output := make([]string, value.Len())
		for i := 0; i < value.Len(); i++ {
			output[i] = fmt.Sprintf("%s", value.Index(i))
		}

		return output
	} else if value.Kind() == reflect.String {
		fmt.Printf("STRING `%v`\n", value.String())
		return []string{value.String()}
	}

	return []string{}
}

func NewJob(filename string, payload metadata.Payload) *Job {
	return &Job{filename, payload}
}
