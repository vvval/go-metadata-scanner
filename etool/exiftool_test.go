package etool

import (
	"github.com/vvval/go-metadata-scanner/util"
	"github.com/vvval/go-metadata-scanner/vars/metadata"
	"reflect"
	"testing"
)

func TestReadArgs(t *testing.T) {
	type check struct {
		names  []string
		fields []string
		exp    []string
	}

	set := []check{
		{[]string{"a", "b"}, []string{"f1", "f2"}, []string{"-j", "-G", "-b", "-f1:all", "-f2:all", "a", "b"}},
		{[]string{}, []string{"f1", "f2"}, []string{"-j", "-G", "-b", "-f1:all", "-f2:all"}},
		{[]string{"a", "b"}, []string{}, []string{"-j", "-G", "-b", "a", "b"}},
	}

	for i, v := range set {
		p := packReadArgs(v.names, v.fields)
		if !reflect.DeepEqual(p, v.exp) {
			t.Errorf("Args mismatch for %d:\ninput %s and %s:\ngot %s\nexpected %s", i, v.names, v.fields, p, v.exp)
		}
	}
}

func TestWriteArgs(t *testing.T) {
	type check struct {
		name      string
		tags      metadata.Tags
		useSep    bool
		originals bool
		exp       []string
	}

	set := []check{
		{"name", metadata.Tags{"n1": "v1", "n2": "v2,v3"}, true, true,
			[]string{"-n1=v1", "-n2=v2,v3", "-sep", metadata.Separator(), "name"}},
		{"", metadata.Tags{"n1": "v1", "n2": "v2,v3"}, true, true,
			[]string{"-n1=v1", "-n2=v2,v3", "-sep", metadata.Separator(), ""}},
		{"", metadata.Tags{}, true, true,
			[]string{"-sep", metadata.Separator(), ""}},
		{"name", metadata.Tags{"n1": "v1", "n2": "v2,v3"}, false, true,
			[]string{"-n1=v1", "-n2=v2,v3", "name"}},
		{"name", metadata.Tags{"n1": "v1", "n2": "v2,v3"}, false, false,
			[]string{"-n1=v1", "-n2=v2,v3", overwriteFlag, "name"}},
	}
	for i, v := range set {
		p := packWriteArgs(v.name, v.tags, v.useSep, v.originals)
		if !util.Equal(p, v.exp) {
			t.Errorf("Args mismatch for %d check\ngot %+v\nexp %+v", i, p, v.exp)
		}
	}
}
