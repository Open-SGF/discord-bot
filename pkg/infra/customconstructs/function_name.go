package customconstructs

type ResourceNamer struct {
	prefix string
	name   string
}

func NewResourceNamer(prefix, name string) *ResourceNamer {
	if prefix != "" {
		prefix += "-"
	}

	return &ResourceNamer{
		prefix: prefix,
		name:   name,
	}
}

func (f *ResourceNamer) Name() *string {
	name := f.name
	return &name
}

func (f *ResourceNamer) FullName() *string {
	name := f.prefix + f.name
	return &name
}
