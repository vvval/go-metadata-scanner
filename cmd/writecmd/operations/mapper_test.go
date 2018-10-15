package operations

import (
	"fmt"
	"github.com/vvval/go-metadata-scanner/config"
	"github.com/vvval/go-metadata-scanner/configuration"
	"github.com/vvval/go-metadata-scanner/vars"
	"github.com/vvval/go-metadata-scanner/vars/metadata"
	"reflect"
	"testing"
)

//broken test
func TestMapPayload(t *testing.T) {
	dict := configuration.Load(config.DictConfig{}, "./../../../dict.yaml").(config.DictConfig)
	columns := readColumns(map[int]string{0: "", 1: "keywords", 2: "", 3: "title", 4: "test", 5: "XMP:Marked"}, dict)

	type check struct {
		data map[int]string
		has  map[string]interface{}
		miss []string
	}

	set := []check{
		{
			map[int]string{0: "name1", 1: "keyword1,keyword2,keyword3", 2: "empty1", 3: "title1", 4: "test1"},
			map[string]interface{}{"IPTC:Keywords": fmt.Sprintf("keyword1%skeyword2%skeyword3", metadata.Separator(), metadata.Separator()), "IPTC:Headline": "title1", "XMP:Marked": false},
			[]string{"test", ""},
		},
		{
			map[int]string{0: "name2", 1: "keyword4", 2: "empty2", 3: "", 4: "", 5: "true"},
			map[string]interface{}{"IPTC:Keywords": "keyword4", "XMP:Marked": true},
			[]string{"IPTC:Headline"},
		},
	}

	for i, s := range set {
		payload := mapPayload(columns, s.data, dict)
		tags := payload.Tags()
		for name, val := range s.has {
			if v, ok := tags.Tag(name); !ok {
				t.Errorf("payload not found (line `%d`):\nexp `%s` `%v`", i, name, val)
			} else if v != val {
				t.Errorf("payload mismatch (line `%d`) for `%s`:\ngot `%v`\nexp `%v`", i, name, v, val)
			}
		}
	}
}

func TestReadColumns(t *testing.T) {
	type tag struct {
		pos   int
		name  string
		found bool
	}
	type check struct {
		cols map[int]string
		exp  []tag
	}

	set := []check{
		{map[int]string{0: "abc", 1: "Keywords"}, []tag{
			{1, "keywords", true},
		}},
		{map[int]string{0: "", 1: "keywords "}, []tag{
			{1, "IPTC:Keywords", true},
		}},
		{map[int]string{0: "", 1: "XMP:Description"}, []tag{
			{1, "description", true},
		}},
		{map[int]string{0: "", 1: "keywords", 2: "", 3: "description", 4: "", 5: ""}, []tag{
			{1, "keywords", true},
			{3, "description", true},
		}},
		{map[int]string{0: "", 1: "keywords", 2: "test"}, []tag{
			{1, "keywords", true},
			{2, "test", false},
		}},
	}

	dict := configuration.Load(config.DictConfig{}, "./../../../dict.yaml").(config.DictConfig)

	for i, s := range set {
		read := readColumns(s.cols, dict)
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
