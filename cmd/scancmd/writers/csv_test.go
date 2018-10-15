package writers

import (
	"encoding/json"
	"github.com/vvval/go-metadata-scanner/vars"
	"github.com/vvval/go-metadata-scanner/vars/metadata"
	"reflect"
	"testing"
)

func TestPackLine(t *testing.T) {
	type check struct {
		f      vars.File
		groups []string
		exp    []map[string]string
		ok     bool
	}

	set := []check{
		{
			vars.NewFile("file1.jpg", metadata.Tags{"test:tag": "some test tag", "XMP:tag1": "xmp tag1", "XMP:tag2": "xmp tag2", "iptc:tag3": "iptc tag3"}),
			[]string{"filename", "XMP", "iptc"},
			[]map[string]string{
				{"tag1": "xmp tag1", "tag2": "xmp tag2"},
				{"tag3": "iptc tag3"},
			},
			true,
		},
		{
			vars.NewFile("file2.png", metadata.Tags{"test:tag": "some test tag", "XMP:tag1": "xmp tag1", "XMP:tag2": "xmp tag2", "iptc:tag2": "iptc tag2"}),
			[]string{"filename", "XMP", "iptc"},
			[]map[string]string{
				{"tag1": "xmp tag1"},
				{"tag2": "iptc tag2"},
			},
			false,
		},
	}

	//Case sensitive comparison
	for i, s := range set {
		p := packCSVLine(&s.f, s.groups)
		pp := convert(p[1:])
		if !mapEqual(pp, s.exp) && s.ok {
			t.Errorf("grouping tags failed (line `%d`) (wrong inequality):\ngot `%v`\nexp `%v` (`%t` `%t`)", i, pp, s.exp, s.ok, mapEqual(pp[1:], s.exp))
		} else if reflect.DeepEqual(pp[1:], s.exp) && !s.ok {
			t.Errorf("grouping tags failed (line `%d`) (wrong equality):\ngot `%v`\nexp `%v` (`%t`)", i, pp, s.exp, s.ok)
		}
	}
}

func convert(i []string) []map[string]string {
	var o []map[string]string
	for _, v := range i {
		m := make(map[string]string)
		json.Unmarshal([]byte(v), &m)
		o = append(o, m)
	}

	return o
}

func mapEqual(a, b []map[string]string) bool {
	if a == nil && b != nil || a != nil && b == nil {
		return false
	}

	if len(a) != len(b) {
		return false
	}

	for _, va := range a {
		if !inMap(va, b) {
			return false
		}
	}

	for _, vb := range b {
		if !inMap(vb, a) {
			return false
		}
	}

	return true
}

func inMap(el map[string]string, m []map[string]string) bool {
	for _, v := range m {
		if reflect.DeepEqual(v, el) {
			return true
		}
	}

	return false
}
