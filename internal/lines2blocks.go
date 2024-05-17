package internal

import (
	"regexp"
	"strings"
)

var blockDetect = regexp.MustCompile(
	// Must start with this
	"^`{3,}" +
		`(` +
		// Is starting block; must contain language and name
		`\s*(\w+)\s*(.+)$` +
		`|` +
		// Is ending block; must not contain anything
		`\s*$` +
		`)`,
)

func ToBlocks(lines []string) (Blocks, error) {
	result := Blocks{}

	codeMode := false

	lineBuffer := make([]string, 0)
	var prevMatch []string

	for _, v := range lines {
		if match := blockDetect.FindStringSubmatch(v); match != nil {
			if codeMode {
				block := CodeBlock{
					Language: prevMatch[2],
					Name:     prevMatch[3],
					Content:  strings.Join(lineBuffer, "\n"),
				}
				result = append(result, block)
			} else {
				block := ProseBlock{
					Content: strings.Join(lineBuffer, "\n"),
				}
				result = append(result, block)
			}
			lineBuffer = make([]string, 0)
			codeMode = !codeMode
			prevMatch = match
		} else {
			lineBuffer = append(lineBuffer, v)
		}
	}

	// Add the final block, if it is a prose block
	if !codeMode {
		result = append(result, ProseBlock{
			Content: strings.Join(lineBuffer, "\n"),
		})
	}
	return result, nil
}
