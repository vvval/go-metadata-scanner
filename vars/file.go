package vars

import (
	"github.com/vvval/go-metadata-scanner/vars/metadata"
	"path/filepath"
)

type File struct {
	filename string //parsed with filepath.ToSlashes
	relPath  string
	tags     metadata.Tags
}

func NewFile(filename string, tags metadata.Tags) File {
	return File{
		filename: filepath.ToSlash(filename),
		tags:     tags,
	}
}

func (f *File) Filename() string {
	return f.filename
}

func (f *File) WithRelPath(base string) {
	rel, err := filepath.Rel(base, f.filename)
	if err != nil {
		f.relPath = filepath.ToSlash(f.filename)
	}

	f.relPath = filepath.ToSlash(rel)
}

func (f *File) RelPath() string {
	return f.relPath
}

func (f *File) Tags() metadata.Tags {
	return f.tags
}
