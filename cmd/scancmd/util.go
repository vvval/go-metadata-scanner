package scancmd

import (
	"encoding/json"
	"github.com/vvval/go-metadata-scanner/vars"
	"github.com/vvval/go-metadata-scanner/vars/metadata"
)

const sourceFileField string = "SourceFile"

func parse(data []byte) []vars.File {
	files := []vars.File{}
	schema := []metadata.Tags{}

	if err := json.Unmarshal(data, &schema); err == nil {
		for _, element := range schema {
			if sf, ok := element[sourceFileField]; ok {
				if filename, ok := sourceFile(sf); ok {
					delete(element, sourceFileField)
					file := vars.NewFile(filename, element)
					files = append(files, file)
				}
			}
		}
	}

	return files
}

func sourceFile(s interface{}) (string, bool) {
	if filename, ok := s.(string); ok {
		return filename, true
	}

	return "", false
}
