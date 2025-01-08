package cmd

/*
Prompt:
- Use github.com/spf13/cobra
- Create an external command var called CleanCmd.
- By default this command takes one or more filenames. This is indicated when running --help as [FILES].
  For each filename provided this way, a function called cleanFile() is called with the filename as argument.
- Alternatively it takes one argument: --all (alias -a), in which case no filenames are allowed.
  In this case It scans the current directory and its subdirectories for files matching the "*.go" filename pattern and runs cleanFile() against each.
- If no arguments are provided, print help for the argument.

- cleanFile() function performs the following actions:
  1. Find the FIRST comment block in the file that has the string "Prompt:" as its first contents (ignore whitelines).
  2. Find the end of that comment block.
  3. Delete anything past that comment block.
  4. Overwrite the original file with the new contents.
*/

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
)

// CleanCmd represents the clean command
var CleanCmd = &cobra.Command{
	Use:   "clean [FILES]",
	Short: "Clean specified .go files or all .go files in the current directory and subdirectories",
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
					cleanFile(path)
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
				cleanFile(filename)
			}
		}
	},
}

func init() {
	CleanCmd.Flags().BoolP("all", "a", false, "Clean all .go files in the current directory and subdirectories")
}

func cleanFile(path string) {
	file, err := os.Open(path)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	inPromptBlock := false
	for scanner.Scan() {
		line := scanner.Text()
		if strings.Contains(line, "Prompt:") {
			inPromptBlock = true
		}
		if inPromptBlock && strings.TrimSpace(line) == "*/" {
			lines = append(lines, line)
			break
		}
		lines = append(lines, line)
	}

	if err := scanner.Err(); err != nil {
		fmt.Println("Error reading file:", err)
		return
	}

	err = os.WriteFile(path, []byte(strings.Join(lines, "\n")), 0644)
	if err != nil {
		fmt.Println("Error writing file:", err)
	}
}
