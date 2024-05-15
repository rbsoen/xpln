package cmd

import (
	"github.com/spf13/cobra"
)

var tangleCmd = &cobra.Command{
	Use:   "tangle",
	Short: "Extract source code",
	Run: func(cmd *cobra.Command, args []string) {
	},
}

func init() {
	rootCmd.AddCommand(tangleCmd)
}
