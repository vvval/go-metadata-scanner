package scancmd

import (
	"encoding/json"
	"github.com/vvval/go-metadata-scanner/metadata"
)

const sourceFileField string = "SourceFile"

func parse(data []byte, chunk Chunk) []FileData {
	files := make([]FileData, len(chunk))
	schema := make([]metadata.Tags, len(chunk))
	if err := json.Unmarshal(data, &schema); err == nil {
		for i, elem := range schema {
			if sf, ok := elem[sourceFileField]; ok {
				delete(elem, sourceFileField)
				file := FileData{filename: sf.(string), tags: elem}
				files[i] = file
			}
		}
	}

	return files
}
