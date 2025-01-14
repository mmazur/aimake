/*
Package "cmd".

Uses github.com/spf13/cobra

Variables:
VersionCmd is a &cobra.Command that implements the "version" cli command.
This command takes no parameters.
Upon execution it prints "0.1.0" on the console.
*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var VersionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("0.1.0")
	},
}
