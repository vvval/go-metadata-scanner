package writers

import (
	"encoding/json"
	"fmt"
	"github.com/vvval/go-metadata-scanner/vars"
	"github.com/vvval/go-metadata-scanner/vars/metadata"
	"reflect"
	"testing"
)

func TestPackLine(t *testing.T) {
	type check struct {
		f      vars.File
		groups []string
		exp    map[string]map[string]string
		ok     bool
	}

	set := []check{
		{
			vars.NewFile("file1.jpg", metadata.Tags{"test:tag": "some test tag", "XMP:tag1": "xmp tag1", "XMP:tag2": "xmp tag2", "iptc:tag3": "iptc tag3"}),
			[]string{"filename", "XMP", "iptc"},
			map[string]map[string]string{"XMP": {"tag1": "xmp tag1", "tag2": "xmp tag2"}, "iptc": {"tag3": "iptc tag3"}},
			true,
		},
		{
			vars.NewFile("file2.png", metadata.Tags{"test:tag": "some test tag", "XMP:tag1": "xmp tag1", "XMP:tag2": "xmp tag2", "iptc:tag2": "iptc tag2"}),
			[]string{"filename", "XMP", "iptc"},
			map[string]map[string]string{"xmp": {"tag1": "xmp tag1"}, "IPTC": {"tag2": "iptc tag2"}},
			false,
		},
	}

	//Case sensitive comparison
	for i, s := range set {
		p := packCSVLine(&s.f, s.groups)
		fmt.Printf("ppp %+v\n", p)
		exp := convertMap(s.exp)
		if !reflect.DeepEqual(p[1:], exp) && s.ok {
			t.Errorf("grouping tags failed (line `%d`) (wrong inequality):\ngot `%v`\nexpected `%t` `%v`", i, p, s.ok, exp)
		} else if reflect.DeepEqual(p[1:], exp) && !s.ok {
			t.Errorf("grouping tags failed (line `%d`) (wrong equality):\ngot `%v`\nexpected `%t` `%v`", i, p, s.ok, exp)
		}
	}
}

func convertMap(i map[string]map[string]string) []string {
	var o []string
	for _, v := range i {
		s, _ := json.Marshal(v)
		o = append(o, string(s))
	}

	return o
}
