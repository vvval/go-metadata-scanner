package operations

import (
	"github.com/vvval/go-metadata-scanner/config"
	"github.com/vvval/go-metadata-scanner/configuration"
	"github.com/vvval/go-metadata-scanner/vars"
	"reflect"
	"testing"
)

func testMapPayload(t *testing.T) {

}

func TestReadColumns(t *testing.T) {
	type tag struct {
		pos   int
		name  string
		found bool
	}
	type check struct {
		cols []string
		exp  []tag
	}

	set := []check{
		{[]string{"abc", "Keywords"}, []tag{
			{1, "keywords", true},
		}},
		{[]string{"", "keywords "}, []tag{
			{1, "IPTC:Keywords", true},
		}},
		{[]string{"", "XMP:Description"}, []tag{
			{1, "description", true},
		}},
		{[]string{"", "keywords", "", "description", "", ""}, []tag{
			{1, "keywords", true},
			{3, "description", true},
		}},
		{[]string{"", "keywords", "test"}, []tag{
			{1, "keywords", true},
			{2, "test", false},
		}},
	}

	dict := configuration.Load(config.Dict, "./../../../dict.yaml").(config.DictConfig)

	for i, s := range set {
		read := readColumns(s.cols)
		for j, checkTag := range s.exp {
			foundCheckTag, found := dict.Find(checkTag.name)
			if found != checkTag.found {
				t.Errorf("columns mismatch (line `%d`, column `%d`):\ngot `%t` `%+v`\nexp `%t` `%+v`", i, j, found, foundCheckTag, checkTag.found, checkTag)
			} else if found {
				mapTag, ok := read[checkTag.pos]
				if !ok {
					t.Errorf("columns not found (line `%d`, column `%d`):\nexp `%+v`", i, j, foundCheckTag)
				} else if !tagsEqual(foundCheckTag, mapTag) {
					t.Errorf("column tags not equal (line `%d`, column `%d`):\ngot `%+v`\nexp `%+v`", i, j, mapTag, foundCheckTag)
				}
			}
		}
	}
}

func tagsEqual(t1, t2 vars.Tag) bool {
	if t1.Key() != t2.Key() {
		return false
	}

	if !reflect.DeepEqual(t1.Map(), t2.Map()) {
		return false
	}

	return true
}
