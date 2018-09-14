package vars

type Tag struct {
	key, original string
	tags          []string
}

func (t Tag) Map() []string {
	return t.tags
}

func (t Tag) Key() string {
	return t.key
}

func (t Tag) Original() string {
	return t.original
}

func NewFoundTag(key, original string, tags []string) Tag {
	return Tag{key, original, tags}
}

func NewNotFoundTag(original string) Tag {
	return Tag{original: original}
}
