package internal

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strings"
	"bytes"

	"github.com/gomarkdown/markdown"
	"github.com/gomarkdown/markdown/html"
	"github.com/gomarkdown/markdown/parser"
)

// To link to other code blocks, use the
// @{Name of code block} syntax.
var linkDetect = regexp.MustCompile(
	`^(\s*)@{(.+)}\s*$`,
)

// These templates may be used in prose blocks
// to insert a special list

// `{{ table of contents }} inserts the ToC from
// all prose blocks
var proseBlockTOCCommand = regexp.MustCompile(
	`(<p>\s*)?\{\{\s*table of contents\s*\}\}(\s*</p>)?`,
)

// `{{ index }} inserts the index of all code
// blocks
var proseBlockIndexCommand = regexp.MustCompile(
	`(<p>\s*)?\{\{\s*index\s*\}\}(\s*<p>)?`,
)

// Used for the "used by" indicators for
// each code block
var codeUseMap = make(map[string][]string)

func Weave(filenames ...string) error {
	lines := make([]string, 0)

	// 1. `cat` every text file together
	// into one giant text document
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

	// 2. Transform the text document into blocks
	// of prose and code
	b, err := ToBlocks(lines)
	if err != nil {
		return err
	}

	index := make([]string, 0)
	wholeProse := ""

	// 3. Track where each code block is used
	// in other code blocks
	// Also, track the names of every code block for
	// the index list and concatenate every prose block
	// to generate a ToC.
	for _, v := range b {
		switch t := v.(type) {
		case ProseBlock:
			wholeProse += t.Content + "\n"
		case CodeBlock:
			index = append(index, t.Name)
			t.traceUsages()
		}
	}

	// 4. Render the results
	for _, v := range b {
		switch t := v.(type) {
		default:
			return fmt.Errorf("unexpected type %T", t)
		case ProseBlock:
			r := t.ToHTML()
			r = proseBlockTOCCommand.ReplaceAllStringFunc(r, func(_ string) string {
				return toTOC(wholeProse)
			})
			r = proseBlockIndexCommand.ReplaceAllStringFunc(r, func(_ string) string {
				x := `<div class="index"><nav><ul>`+"\n"
				for _, v := range index {
					x += `<li><a href="#`+makeIdentifier(v)+`">`+v+`</a></li>`+"\n"
				}
				x += `</ul></nav></div>`
				return  x
			})
			fmt.Println(r)
		case CodeBlock:
			fmt.Println(t.ToHTML())
		}
	}
	return err
}

func (b CodeBlock) traceUsages() {
	// Match on a per-line basis
	k := strings.Split(b.Content, "\n")

	for i := 0; i < len(k); i++ {
		if match := linkDetect.FindStringSubmatch(k[i]); match != nil {
			// Update the key of the *referenced block* to add *this block*
			codeUseMap[match[2]] = append(codeUseMap[match[2]], b.Name)
		}
	}
}

func (b CodeBlock) ToHTML() string {
	result := make([]string, 0)
	result = append(result, `<div class="codeblock" id="`+makeIdentifier(b.Name)+`">`)

	result = append(result, `<header class="codeblock-title">`)
	result = append(result, `<a href="#`+makeIdentifier(b.Name)+`">`+b.Name+`</a>`)
	result = append(result, `</header>`)
	result = append(result, `<pre><code class="language-`+b.Language+`">`)

	// Make the resulting code display as-is in the HTML
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

	result = append(result, `<footer class="codeblock-footer">`)
	// If this code block is referenced in other code blocks, add
	// a "Used by" footer.
	if n := len(codeUseMap[b.Name]); n > 0 {
		result = append(result, `<span>Used by </span><ul>`)
		for i := 0; i < n; i++ {
			which := codeUseMap[b.Name][i]
			result = append(result, fmt.Sprintf(`<li><a href="#%s">%s</a></li>`, makeIdentifier(which), which))
		}
		result = append(result, `</ul>`)
	}
	result = append(result, `</footer>`)

	result = append(result, `</div>`)
	return strings.Join(result, "\n")
}

// Turn the string into an HTML-safe identifier
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

// Converting a prose block to HTML is relatively simple
func (b ProseBlock) ToHTML() string {
	p := parser.NewWithExtensions(parser.CommonExtensions | parser.AutoHeadingIDs)
	doc := p.Parse([]byte(b.Content))
	r := html.NewRenderer(
		html.RendererOptions{Flags: html.CommonFlags},
	)
	return string(markdown.Render(doc, r))
}

func toTOC (text string) string {
	p := parser.NewWithExtensions(parser.CommonExtensions | parser.AutoHeadingIDs)
	doc := p.Parse([]byte(text))
	r := html.NewRenderer(
		html.RendererOptions{Flags: html.TOC},
	)
	var b bytes.Buffer
	r.RenderHeader(&b, doc)
	return `<div class="toc">` + b.String() + `</div>`
}
