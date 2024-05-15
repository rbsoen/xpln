package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var tangleCmd = &cobra.Command{
	Use:   "tangle",
	Short: "Extract source code",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("not implemented yet")
	},
}

func init() {
	rootCmd.AddCommand(tangleCmd)
}
