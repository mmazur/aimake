package cmd

/*
Prompt:
- Use github.com/spf13/cobra
- Create an external command var called VersionCmd.
- This command takes no parameters, and prints "0.1.0" on the console.
*/

import (
	"fmt"

	"github.com/spf13/cobra"
)

// VersionCmd represents the version command
var VersionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number of aimake",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("0.1.0")
	},
}
