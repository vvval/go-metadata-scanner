package dict

type tag struct {
	key, original string
	tags          []string
}

func (t tag) Map() []string {
	return t.tags
}

func (t tag) Key() string {
	return t.key
}

func (t tag) Original() string {
	return t.original
}
