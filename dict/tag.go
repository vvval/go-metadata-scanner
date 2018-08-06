package dict

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
