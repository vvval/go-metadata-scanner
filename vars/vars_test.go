package vars

import (
	"github.com/vvval/go-metadata-scanner/vars/metadata"
	"reflect"
	"testing"
)

func TestChunk(t *testing.T) {
	chunk := Chunk{"a", "b", "c", "d"}

	setEqual := [][][]Chunk{
		{chunk.Split(1), {Chunk{"a"}, Chunk{"b"}, Chunk{"c"}, Chunk{"d"}}},
		{chunk.Split(2), {Chunk{"a", "b"}, Chunk{"c", "d"}}},
		{chunk.Split(3), {Chunk{"a", "b", "c"}, Chunk{"d"}}},
		{chunk.Split(4), {Chunk{"a", "b", "c", "d"}}},
		{chunk.Split(5), {Chunk{"a", "b", "c", "d"}}},
	}

	for _, s := range setEqual {
		if !reflect.DeepEqual(s[0], s[1]) {
			t.Errorf("chunks should be equal:\ngot `%s`\nexpected `%s`", s[0], s[1])
		}
	}

	setUnequal := [][][]Chunk{
		{chunk.Split(1), {Chunk{"a"}, Chunk{"b"}, Chunk{"c", "d"}}},
		{chunk.Split(2), {Chunk{"a", "b"}, Chunk{"c", "d"}, Chunk{}}},
		{chunk.Split(3), {Chunk{"a", "b", "d"}, Chunk{"e"}}},
	}

	for _, s := range setUnequal {
		if reflect.DeepEqual(s[0], s[1]) {
			t.Errorf("chunks should not be equal:\ngot `%s`\nexpected `%s`", s[0], s[1])
		}
	}
}

func TestTag(t *testing.T) {
	tags := metadata.Tags{"a": "b"}
	v, ok := tags.Tag("a")
	if !ok {
		t.Errorf("tag `%s` not found\nshould be `%s`", "a", "b")
	}

	if v != "b" {
		t.Errorf("tag `%s` found but value is wrong:\ngot `%s`\nexpected `%s`", "a", v, "b")
	}
}
