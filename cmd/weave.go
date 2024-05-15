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
		err := internal.Weave()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(weaveCmd)
}
