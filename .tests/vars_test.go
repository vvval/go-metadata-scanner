package _tests

import (
	"github.com/vvval/go-metadata-scanner/vars"
	"github.com/vvval/go-metadata-scanner/vars/metadata"
	"reflect"
	"testing"
)

func TestFile(t *testing.T) {
	filename := "file\\with/different/slashes\\usage"
	file := vars.NewFile(filename, metadata.Tags{})
	filenameSlashes := "file/with/different/slashes/usage"

	if file.Filename() != filenameSlashes {
		t.Errorf("filename should be converted to slashes format:\ngot `%s`\nexpected `%s`", file.Filename(), filenameSlashes)
	}

	dir := "file/with"
	file.WithRelPath(dir)
	relPath := "different/slashes/usage"

	if file.RelPath() != relPath {
		t.Errorf("rel path should converted to slashes format:\ngot `%s`\nexpected `%s`", file.RelPath(), relPath)
	}
}

func TestChunk(t *testing.T) {
	chunk := vars.Chunk{"a", "b", "c", "d"}

	setEqual := [][][]vars.Chunk{
		{chunk.Split(1), {vars.Chunk{"a"}, vars.Chunk{"b"}, vars.Chunk{"c"}, vars.Chunk{"d"}}},
		{chunk.Split(2), {vars.Chunk{"a", "b"}, vars.Chunk{"c", "d"}}},
		{chunk.Split(3), {vars.Chunk{"a", "b", "c"}, vars.Chunk{"d"}}},
		{chunk.Split(4), {vars.Chunk{"a", "b", "c", "d"}}},
		{chunk.Split(5), {vars.Chunk{"a", "b", "c", "d"}}},
	}

	for _, s := range setEqual {
		if !reflect.DeepEqual(s[0], s[1]) {
			t.Errorf("chunks should be equal:\ngot `%s`\nexpected `%s`", s[0], s[1])
		}
	}

	setUnequal := [][][]vars.Chunk{
		{chunk.Split(1), {vars.Chunk{"a"}, vars.Chunk{"b"}, vars.Chunk{"c", "d"}}},
		{chunk.Split(2), {vars.Chunk{"a", "b"}, vars.Chunk{"c", "d"}, vars.Chunk{}}},
		{chunk.Split(3), {vars.Chunk{"a", "b", "d"}, vars.Chunk{"e"}}},
	}

	for _, s := range setUnequal {
		if reflect.DeepEqual(s[0], s[1]) {
			t.Errorf("chunks should not be equal:\ngot `%s`\nexpected `%s`", s[0], s[1])
		}
	}
}

func TestMetadataTag(t *testing.T) {
	tags := metadata.Tags{"a": "b"}
	v, ok := tags.Tag("a")
	if !ok {
		t.Errorf("tag `%s` not found\nshould be `%s`", "a", "b")
	}

	if v != "b" {
		t.Errorf("tag `%s` found but value is wrong:\ngot `%s`\nexpected `%s`", "a", v, "b")
	}
}

func TestPayloadAddTag(t *testing.T) {
	p := metadata.New()

	p.AddTag("tag1", "true")
	if p.Tags()["tag1"] != true {
		t.Errorf("`true` strings should be converted to bool type, got %b", p.Tags()["tag1"])
	}

	p.AddTag("tag2", "false")
	if p.Tags()["tag2"] != false {
		t.Errorf("`false` strings should be converted to bool type, got %b", p.Tags()["tag2"])
	}

	p.AddTag("tag3", "test")
	if p.Tags()["tag3"] != "test" {
		t.Errorf("regular strings should not be converted, got %s", p.Tags()["tag3"])
	}
}

func TestPayloadAddList(t *testing.T) {
	p := metadata.New()
	if p.UseSeparator() {
		t.Error("separator usage should be `false` by default")
	}

	p.AddList("tag1", []string{"a"})
	if p.UseSeparator() {
		t.Error("separator usage should not be set to `true` if adding a list tag with less than 2 non-empty values")
	}

	p.AddList("tag2", []string{"a", "a", "", " "})
	if p.UseSeparator() {
		t.Errorf("separator usage should not be set to `true` if adding a list tag with less than 2 non-empty unique values: %s", p.Tags()["tag2"])
	}

	p.AddList("tag3", []string{"a", "b"})
	if !p.UseSeparator() {
		t.Error("separator should be set to `true` after adding a list tag with more than one non-empty value")
	}

	p.AddList("tag4", []string{"c"})
	if !p.UseSeparator() {
		t.Error("separator shouldn't be set to `false` once it is set to `true`")
	}
}

func TestPayloadUpdateTag(t *testing.T) {
	p := metadata.New()
	if p.UseSeparator() {
		t.Error("separator usage should be `false` by default")
	}

	p.AddList("tag1", []string{"a"})
	if p.UseSeparator() {
		t.Error("separator usage should not be set to `true` if adding a list tag with less than 2 non-empty values")
	}

	p.UpdateList("tag1", []string{"b"})
	if !p.UseSeparator() {
		t.Error("separator should be set to `true` after adding a list tag with more than one non-empty value")
	}

	p.AddList("tag1", []string{"c"})
	if p.UseSeparator() {
		t.Error("separator usage should not be set to `true` if adding a list tag with less than 2 non-empty values")
	}
}
