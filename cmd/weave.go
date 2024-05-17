package cmd

import (
	"fmt"
	"os"
	"xpln/internal"

	"github.com/spf13/cobra"
)

var weaveCmd = &cobra.Command{
	Use:   "weave source1.md source2.md...",
	Short: "Generate a web page",
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
		fmt.Println(`
		<!DOCTYPE html>
		<html>
			<head>
				<meta charset="utf-8">
				<meta name="viewport" content="width=device-width,height=device-height,initial-scale=1">
			</head>
			<body>
		`)
		fmt.Println(output)
		fmt.Println(`
			</body>
		</html>
		`)
	},
}

func init() {
	rootCmd.AddCommand(weaveCmd)
}
