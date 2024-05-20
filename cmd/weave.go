package cmd

import (
	"fmt"
	"os"
	"xpln/internal"

	"github.com/spf13/cobra"
	"github.com/yosssi/gohtml"
)

var weaveCmd = &cobra.Command{
	Use:   "weave source1.md source2.md...",
	Short: "Generate a web page",
	Long:`
Generates an HTML document out of a literate program.

This command does not generate separate HTML documents, but rather concatenates
all parts into a single document for reading on the web or for printing into
a PDF.

A minimal shell will be provided for the document. It assumes, in the same
directory, a style.css file for use on the web, as well as a style.print.css
for printing.`,
	Example: "xpln weave 001_intro.md 002_explanation.md 003_conclusion > index.html",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 1 {
			fmt.Println("must provide at least 1 file")
			os.Exit(1)
		}
		output, err := internal.Weave(args...)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		// TODO: configurable template
		out := ""
		out += `
		<!DOCTYPE html>
		<html>
			<head>
				<meta charset="utf-8">
				<meta name="viewport" content="width=device-width,height=device-height,initial-scale=1">
				<link rel="stylesheet" href="style.css" media="screen">
				<link rel="stylesheet" href="style.print.css" media="print">
			</head>
			<body>
		`
		out += output
		out += `
			</body>
		</html>
		`
		fmt.Println(gohtml.Format(out))
	},
}

func init() {
	rootCmd.AddCommand(weaveCmd)
}
