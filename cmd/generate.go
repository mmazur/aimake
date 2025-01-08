package cmd

/*
Prompt:
- Use github.com/spf13/cobra
- Create an external command var called GenerateCmd.
- Full name of the command is 'generate', but it also is available as an alias 'gen'.
- By default this command takes one or more filenames. This is indicated when running --help as [FILES].
  For each filename provided this way, a function called generateFile() is called with the filename as argument.
- Alternatively it takes one argument: --all (alias -a), in which case no filenames are allowed.
  In this case It scans the current directory and its subdirectories for files matching the "*.go" filename pattern and runs generateFile() against each.
- If no arguments are provided, print help for the argument.

- generateFile() function performs the following actions:
  Print the filename.
*/

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

// GenerateCmd represents the generate command
var GenerateCmd = &cobra.Command{
	Use:     "generate [FILES]",
	Aliases: []string{"gen"},
	Short:   "Generate specified .go files or all .go files in the current directory and subdirectories",
	Run: func(cmd *cobra.Command, args []string) {
		all, _ := cmd.Flags().GetBool("all")
		if all {
			if len(args) > 0 {
				fmt.Println("Error: --all flag cannot be used with filenames")
				return
			}
			err := filepath.Walk(".", func(path string, info os.FileInfo, err error) error {
				if err != nil {
					return err
				}
				if filepath.Ext(path) == ".go" {
					generateFile(path)
				}
				return nil
			})
			if err != nil {
				fmt.Println("Error:", err)
			}
		} else if len(args) == 0 {
			cmd.Help()
		} else {
			for _, filename := range args {
				generateFile(filename)
			}
		}
	},
}

func init() {
	GenerateCmd.Flags().BoolP("all", "a", false, "Generate all .go files in the current directory and subdirectories")
}

func generateFile(path string) {
	fmt.Println(path)
}
