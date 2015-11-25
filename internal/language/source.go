package language

type Source struct {
	Body string
	Name string
}

func NewSource(body string, name string) Source {
	if name == "" {
		name = "GraphQL"
	}
	return Source{
		Body: body,
		Name: name,
	}
}
