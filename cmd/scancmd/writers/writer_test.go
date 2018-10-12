package writers

import (
	"github.com/vvval/go-metadata-scanner/vars"
	"github.com/vvval/go-metadata-scanner/vars/metadata"
	"reflect"
	"testing"
)

func TestTagsByGroups(t *testing.T) {
	type check struct {
		f      vars.File
		groups []string
		exp    map[string]map[string]string
		ok     bool
	}

	set := []check{
		{
			vars.NewFile("", metadata.Tags{"test:tag": "some test tag", "XMP:tag1": "xmp tag1", "iptc:tag2": "iptc tag2"}),
			[]string{"xmp", "iptc"},
			map[string]map[string]string{"XMP": {"tag1": "xmp tag1"}, "iptc": {"tag2": "iptc tag2"}},
			true,
		},
		{
			vars.NewFile("", metadata.Tags{"test:tag": "some test tag", "XMP:tag1": "xmp tag1", "iptc:tag2": "iptc tag2"}),
			[]string{"xmp", "iptc"},
			map[string]map[string]string{"xmp": {"tag1": "xmp tag1"}, "IPTC": {"tag2": "iptc tag2"}},
			false,
		},
	}

	//Case sensitive comparison
	for i, s := range set {
		c := tagsByGroups(&s.f, s.groups)
		if !reflect.DeepEqual(c, s.exp) && s.ok {
			t.Errorf("grouping tags failed (line `%d`) (wrong inequality):\ngot `%v`\nexpected `%t` `%v`", i, c, s.ok, s.exp)
		} else if reflect.DeepEqual(c, s.exp) && !s.ok {
			t.Errorf("grouping tags failed (line `%d`) (wrong equality):\ngot `%v`\nexpected `%t` `%v`", i, c, s.ok, s.exp)
		}
	}
}
