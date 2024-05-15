package internal

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/gomarkdown/markdown"
	"github.com/gomarkdown/markdown/html"
	"github.com/gomarkdown/markdown/parser"
)

var linkDetect = regexp.MustCompile(
	`^(\s*)@{(.+)}\s*$`,
)

func (b ProseBlock) ToHTML() string {
	p := parser.NewWithExtensions(parser.CommonExtensions | parser.AutoHeadingIDs)
	doc := p.Parse([]byte(b.Content))
	r := html.NewRenderer(
		html.RendererOptions{Flags: html.CommonFlags},
	)

	return string(markdown.Render(doc, r))
}

func makeIdentifier(id string) string {
	result := ""
	lower := strings.ToLower(id)
	for i := 0; i < len(lower); i++ {
		if ((lower[i] >= '0') && (lower[i] <= '9')) ||
			((lower[i] >= 'a') && (lower[i] <= 'z')) {
			result = result + string(lower[i])
		} else {
			result = result + "-"
		}
	}
	return result
}

func (b CodeBlock) ToHTML() string {
	result := make([]string, 0)
	result = append(result, `<div class="codeblock" id="`+makeIdentifier(b.Name)+`">`)
	result = append(result, `<header class="codeblock-title">`)
	result = append(result, `<a href="#`+makeIdentifier(b.Name)+`">`+b.Name+`</a>`)
	result = append(result, `</header>`)
	result = append(result, `<pre><code class="language-`+b.Language+`">`)

	s := strings.ReplaceAll(b.Content, `&`, `&amp;`)
	s = strings.ReplaceAll(s, `<`, `&lt;`)
	s = strings.ReplaceAll(s, `>`, `&gt;`)

	k := strings.Split(s, "\n")
	for i := 0; i < len(k); i++ {
		if match := linkDetect.FindStringSubmatch(k[i]); match != nil {
			k[i] = fmt.Sprintf(`%s<a href="#%s">&#12298; %s &#12299;</a>`, match[1], makeIdentifier(match[2]), match[2])
		}
	}

	result = append(result, strings.Join(k, "\n"))
	result = append(result, "</code></pre>")
	result = append(result, `</div>`)
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
			fmt.Println(t.ToHTML())
		case CodeBlock:
			fmt.Println(t.ToHTML())
		}
	}
	return err
}
