package internal

type ProseBlock struct {
	Content string
}

type CodeBlock struct {
	Name     string
	Language string
	Content  string
}

type Renderable interface {
	ToHTML() string
}

type Blocks []interface{}
