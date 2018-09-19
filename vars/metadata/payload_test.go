package metadata

import (
	"testing"
)

func TestPayloadAddTag(t *testing.T) {
	p := New()

	p.AddTag("tag1", "true")
	if p.Tags()["tag1"] != true {
		t.Errorf("`true` strings should be converted to bool type, got %t", p.Tags()["tag1"])
	}

	p.AddTag("tag2", "false")
	if p.Tags()["tag2"] != false {
		t.Errorf("`false` strings should be converted to bool type, got %t", p.Tags()["tag2"])
	}

	p.AddTag("tag3", "test")
	if p.Tags()["tag3"] != "test" {
		t.Errorf("strings should not be converted, got %s", p.Tags()["tag3"])
	}
}

func TestPayloadAddList(t *testing.T) {
	p := New()
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
	p := New()
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
