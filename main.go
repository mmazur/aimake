package main

/*
Prompt:
- This file is the entrypoint to the code and contains the main executable function.
- The app name is "aimake".
- The app is a CLI tool that parses commandline arguments using github.com/spf13/cobra library.
- Supported commands can be found as public functions in the format *Cmd inside the "cmd" package.
- Do not generate any other commands.

Hints:
- Don't forget to define and initiate the root cobra Cmd.
*/

import (
	"fmt"
	"os"

	"aimake/cmd"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "aimake",
	Short: "Make code with AI",
}

func main() {
	rootCmd.AddCommand(cmd.VersionCmd)
	rootCmd.AddCommand(cmd.QueryCmd)
	rootCmd.AddCommand(cmd.CleanCmd)
	rootCmd.AddCommand(cmd.GenerateCmd)

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
