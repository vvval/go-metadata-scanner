package vars

import "github.com/vvval/go-metadata-scanner/vars/metadata"

type File struct {
	filename string
	tags     metadata.Tags
}

func (f *File) Filename() string {
	return f.filename
}

func (f *File) Tags() metadata.Tags {
	return f.tags
}

func NewFile(filename string, tags metadata.Tags) File {
	return File{filename, tags}
}
