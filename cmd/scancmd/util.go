package scancmd

import (
	"encoding/json"
	"github.com/vvval/go-metadata-scanner/vars"
	"github.com/vvval/go-metadata-scanner/vars/metadata"
)

const sourceFileField string = "SourceFile"

func parse(data []byte, chunk vars.Chunk) []vars.File {
	files := make([]vars.File, len(chunk))
	schema := make([]metadata.Tags, len(chunk))

	if err := json.Unmarshal(data, &schema); err == nil {
		for i, element := range schema {
			if sf, ok := element[sourceFileField]; ok {
				delete(element, sourceFileField)
				file := vars.NewFile(sf.(string), element)
				files[i] = file
			}
		}
	}

	return files
}
