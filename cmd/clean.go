package cmd

/*
Prompt:
- Use github.com/spf13/cobra
- Create an external command var called CleanCmd.
- This command takes no arguments.
- It scans the current directory and its subdirectories for files matching the "*.go" filename pattern.
- It prints the paths of all found files.
*/

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

// CleanCmd represents the clean command
var CleanCmd = &cobra.Command{
	Use:   "clean",
	Short: "Scan and print all .go files in the current directory and subdirectories",
	Run: func(cmd *cobra.Command, args []string) {
		err := filepath.Walk(".", func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if filepath.Ext(path) == ".go" {
				fmt.Println(path)
			}
			return nil
		})
		if err != nil {
			fmt.Println("Error:", err)
		}
	},
}

// ...existing code...
