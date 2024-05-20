package cmd

import (
	"fmt"
	"os"
	"xpln/internal"

	"github.com/spf13/cobra"
)

var tangleCmd = &cobra.Command{
	Use:   "tangle out_dir source1.md source2.md...",
	Short: "Extract source code",
	Long: `
Extracts usable source code out of a literate program.

Files are generated out of code blocks that start with a forward slash (/), and
any code references within are evaluated recursively.

Do be careful not to reference its own code block!`,
	Example: "tangle src 001_intro.md 002_explanation.md 003_conclusion.md",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 1 {
			fmt.Println("must provide output directory")
			os.Exit(1)
		}
		if len(args) < 2 {
			fmt.Println("must provide at least 1 file")
			os.Exit(1)
		}
		err := internal.Tangle(args[0], args[1:]...)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(tangleCmd)
}
