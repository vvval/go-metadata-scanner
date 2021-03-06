package operations

import (
	"github.com/vvval/go-metadata-scanner/config"
	"github.com/vvval/go-metadata-scanner/util"
	"github.com/vvval/go-metadata-scanner/vars"
	"github.com/vvval/go-metadata-scanner/vars/metadata"
	"strings"
)

// Script will pick all values for line data that are presented in columns (by column index)
// and fill output array with line value for each tag aliases from TagMap.
//
// Example:
// 		Flags line is [1:"keywords values", 3:"some description"
// 		Flags columns are [1:keywords,2:title,3:description]
// 		Output is [
// 			"IPTC:keywords1":"keywords values",
// 			"XMP:keywords2":"keywords values",
// 			"IPTC:description1":"some description",
// 			"XMP:description2":"some description"
// 		]
func mapPayload(columns map[int]vars.Tag, input map[int]string, dict config.DictConfig) metadata.Payload {
	payload := metadata.New()

	for index, columnTag := range columns {
		value, ok := input[index]
		if !ok {
			value = ""
		}

		for _, tag := range columnTag.Map() {
			if dict.IsBoolean(columnTag.Key(), tag) {
				payload.AddBool(tag, len(value) != 0)
			} else if len(value) != 0 {
				if dict.IsList(columnTag.Key(), tag) {
					//Not AddList because we can have multiple columns for list tags like: "Keywords: Poses", "Keywords: Age"
					payload.UpdateList(tag, util.SplitKeywords(value))
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
func readColumns(columns map[int]string, dict config.DictConfig) map[int]vars.Tag {
	output := map[int]vars.Tag{}
	for i, column := range columns {
		column = strings.Trim(column, " ")
		if i == 0 || len(column) == 0 {
			// Skip first column and empty columns
			continue
		}

		tag, found := dict.Find(column)
		if !found {
			// Skip not found columns
			continue
		}

		output[i] = tag
	}

	return output
}
