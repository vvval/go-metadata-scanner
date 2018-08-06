package writeCommand

import (
	"github.com/vvval/go-metadata-scanner/dict"
	"github.com/vvval/go-metadata-scanner/metadata"
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
func mapPayload(columns map[int]dict.Tag, input []string) metadata.Payload {
	payload := metadata.New()
	d := dict.Get()

	for index, value := range input {
		t, ok := columns[index]
		if !ok {
			// Unmapped key, skip
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
func readColumns(columns []string) map[int]dict.Tag {
	d := dict.Get()
	output := map[int]dict.Tag{}
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

		output[i] = tag
	}

	return output
}
