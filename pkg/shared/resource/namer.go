package resource

type Namer struct {
	prefix string
	name   string
}

func NewNamer(prefix, name string) *Namer {
	if prefix != "" {
		prefix += "-"
	}

	return &Namer{
		prefix: prefix,
		name:   name,
	}
}

func (f *Namer) Name() string {
	return f.name
}

func (f *Namer) FullName() string {
	return f.prefix + f.name
}
