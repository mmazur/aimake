package cmd

/*
Prompt:
- Use github.com/spf13/cobra
- Create an external command var called CleanCmd.
- This command takes no arguments.
- It scans the current directory and its subdirectories for files matching the "*.go" filename pattern.
- For each file found it performs the following actions (implemented as a separate, private function):
  0. Print the path of the file.
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
	Use:   "clean",
	Short: "Scan and clean all .go files in the current directory and subdirectories",
	Run: func(cmd *cobra.Command, args []string) {
		err := filepath.Walk(".", func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if filepath.Ext(path) == ".go" {
				fmt.Println(path)
				cleanFile(path)
			}
			return nil
		})
		if err != nil {
			fmt.Println("Error:", err)
		}
	},
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
