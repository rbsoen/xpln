package internal

import (
	"bufio"
	"os"
	"path/filepath"
	"strings"
)

var codeContentMap = make(map[string][]string)

func Tangle(outputFolder string, filenames ...string) error {
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

	// 3. Track where each code block is used
	// in other code blocks
	for _, v := range b {
		switch t := v.(type) {
		case CodeBlock:
			codeContentMap[t.Name] = strings.Split(t.Content, "\n")
		}
	}

	// 4. Resolve code block references
	for _, v := range b {
		switch t := v.(type) {
		case CodeBlock:
			if t.Name[0] == '/' {
				content := strings.Join(resolveReferences(t.Name, ""), "\n")
				dir, file := filepath.Split(t.Name[1:])
				outDir := filepath.Join(outputFolder, dir)
				if err := os.MkdirAll(outDir, os.ModePerm); err != nil {
					return err
				}
				if err := os.WriteFile(filepath.Join(outDir, file), []byte(content), os.ModePerm); err != nil {
					return err
				}
			}
		}
	}
	return nil
}

func resolveReferences(blockName string, leading string) []string {
	result := make([]string, 0)
	for _, v := range codeContentMap[blockName] {
		if match := linkDetect.FindStringSubmatch(v); match != nil {
			result = append(result, resolveReferences(match[2], match[1])...)
		} else {
			result = append(result, leading+v)
		}
	}
	return result
}
