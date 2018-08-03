package writer

import (
	"github.com/vvval/go-metadata-scanner/cmd/metadata"
	"github.com/vvval/go-metadata-scanner/dict"
	"github.com/vvval/go-metadata-scanner/util"
	"strings"
)

// Script will pick all values for line data that are presented in columns (by column index)
// and fill output array with line value for each tag aliases from TagMap.
//
// Example:
// 		Input line is [1:"keywords values", 3:"some description"
// 		Input columns are [1:keywords,2:title,3:description]
// 		Output is [
// 			"IPTC:keywords1":"keywords values",
// 			"XMP:keywords2":"keywords values",
// 			"IPTC:description1":"some description",
// 			"XMP:description2":"some description"
// 		]
func MapPayload(columns map[int]string, input []string) metadata.Payload {
	payload := metadata.New()
	d := dict.Get()

	for index, value := range input {
		key, ok := columns[index]
		if !ok {
			// Unmapped key, skip
			continue
		}

		t, found := d.Find(key)
		if !found {
			// Unknown tag, skip
			continue
		}

		//todo replace with getLists and getBooleans
		for _, tag := range t.Map() {
			if d.IsBoolean(t.Key(), tag) {
				payload.AddBool(tag, len(value) != 0)
			} else if len(value) != 0 {
				if d.IsList(t.Key(), tag) {
					payload.AddList(tag, util.SplitKeywords(value))
				} else {
					payload.AddTag(tag, value)
				}
			}
		}
	}

	return payload
}

// Map columns to a known tag map
// Skip 1st column (dedicated to a file names) and empty columns
func ReadColumns(columns []string) map[int]string {
	d := dict.Get()
	output := map[int]string{}
	for i, column := range columns {
		column = strings.Trim(column, " ")
		if i == 0 || len(column) == 0 {
			// Skip first column and empty columns
			continue
		}

		tag, found := d.Find(column)
		if !found {
			// Skip not found columns
			continue
		}

		output[i] = tag.Original()
	}

	return output
}
