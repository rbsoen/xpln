package cmd

import (
	"fmt"
	"os"
	"xpln/internal"

	"github.com/spf13/cobra"
)

var weaveCmd = &cobra.Command{
	Use:   "weave",
	Short: "Generate a web page",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 1 {
			fmt.Println("must provide at least 1 file")
			os.Exit(1)
		}
		err := internal.Weave(args...)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(weaveCmd)
}
