package internal

import (
	"fmt"
	"os"

	"github.com/gomarkdown/markdown"
	"github.com/gomarkdown/markdown/html"
	"github.com/gomarkdown/markdown/parser"
)

func Weave() error {
	i, err := os.ReadFile("test.md")
	if err != nil {
		return err
	}

	p := parser.NewWithExtensions(parser.CommonExtensions | parser.AutoHeadingIDs)
	doc := p.Parse(i)
	r := html.NewRenderer(
		html.RendererOptions{Flags: html.CommonFlags},
	)

	fmt.Print(string(markdown.Render(doc, r)))
	return nil
}
