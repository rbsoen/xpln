package internal

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/gomarkdown/markdown"
	"github.com/gomarkdown/markdown/html"
	"github.com/gomarkdown/markdown/parser"
)

func (b ProseBlock) ToHTML() string {
	p := parser.NewWithExtensions(parser.CommonExtensions | parser.AutoHeadingIDs)
	doc := p.Parse([]byte(b.Content))
	r := html.NewRenderer(
		html.RendererOptions{Flags: html.CommonFlags},
	)

	return string(markdown.Render(doc, r))
}

func (b CodeBlock) ToHTML() string {
	result := make([]string, 0)
	result = append(result, `<pre><code class="language-`+b.Language+`">`)
	// TODO: escape characters
	result = append(result, b.Content)
	result = append(result, "</code></pre>")
	return strings.Join(result, "\n")
}

func Weave(filenames ...string) error {
	lines := make([]string, 0)

	for _, filename := range filenames {
		file, err := os.Open(filename)
		if err != nil {
			return err
		}
		defer file.Close()

		s := bufio.NewScanner(file)
		for s.Scan() {
			lines = append(lines, s.Text())
		}
	}

	b, err := ToBlocks(lines)
	if err != nil {
		return err
	}

	for _, v := range b {
		switch t := v.(type) {
		default:
			return fmt.Errorf("unexpected type %T", t)
		case ProseBlock:
			fmt.Print(t.ToHTML())
		case CodeBlock:
			fmt.Print(t.ToHTML())
		}
	}
	return err
}
