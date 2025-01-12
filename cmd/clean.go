/*
Prompt:
  - This is `package cmd`
  - Use github.com/spf13/cobra
  - Create an external command var called CleanCmd.
  - By default this command takes one or more filenames. This is indicated when running --help as [FILES].
    For each filename provided this way, a function called cleanFile() is called with the filename as argument.
  - Alternatively it takes one argument: --all (alias -a), in which case no filenames are allowed.
    In this case It scans the current directory and its subdirectories for files matching the "*.go" filename pattern and runs cleanFile() against each.
  - If no arguments are provided, print help for the argument.

- cleanFile() function performs the following actions:
 1. Find the FIRST comment block in the file that has the string "Prompt:" as its first contents (on the same line or on the next line).
 2. Find the end of that comment block.
 3. Delete anything past that comment block.
 4. Overwrite the original file with the new contents.

Hints:
- Don't bother supporting // style comments.
- There's no such thing as strings.NewScanner
*/
package cmd

import (
	"bufio"
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
)

var CleanCmd = &cobra.Command{
	Use:   "clean [FILES]",
	Short: "Clean Go files by removing content after the 'Prompt:' comment block.",
	Long: `The clean command removes all content from the specified Go files after the first comment block
that starts with 'Prompt:'. You can provide one or more files as arguments or use the --all flag to
clean all .go files in the current directory and its subdirectories.`,
	Args: cobra.ArbitraryArgs,
	Run: func(cmd *cobra.Command, args []string) {
		all, _ := cmd.Flags().GetBool("all")
		if all {
			if len(args) > 0 {
				fmt.Println("When using --all, no filenames should be provided.")
				cmd.Help()
				return
			}
			err := filepath.Walk(".", func(path string, info os.FileInfo, err error) error {
				if err != nil {
					return err
				}
				if !info.IsDir() && strings.HasSuffix(info.Name(), ".go") {
					cleanFile(path)
				}
				return nil
			})
			if err != nil {
				fmt.Printf("Error walking the path: %v\n", err)
			}
		} else if len(args) > 0 {
			for _, file := range args {
				cleanFile(file)
			}
		} else {
			cmd.Help()
		}
	},
}

func init() {
	CleanCmd.Flags().BoolP("all", "a", false, "Clean all .go files in the current directory and subdirectories")
}

func cleanFile(filename string) {
	content, err := ioutil.ReadFile(filename)
	if err != nil {
		fmt.Printf("Error reading file %s: %v\n", filename, err)
		return
	}

	scanner := bufio.NewScanner(bytes.NewReader(content))
	var output bytes.Buffer
	foundPrompt := false
	inCommentBlock := false

	for scanner.Scan() {
		line := scanner.Text()
		if !foundPrompt {
			if strings.Contains(line, "/*") {
				inCommentBlock = true
				commentStartIndex := strings.Index(line, "/*") + 2
				remaining := strings.TrimSpace(line[commentStartIndex:])
				if strings.HasPrefix(remaining, "Prompt:") {
					foundPrompt = true
				} else {
					// Check the next line
					if scanner.Scan() {
						nextLine := strings.TrimSpace(scanner.Text())
						if strings.HasPrefix(nextLine, "Prompt:") {
							foundPrompt = true
						}
						output.WriteString(line + "\n")
						line = nextLine
					}
				}
			}
		}

		if foundPrompt && inCommentBlock {
			output.WriteString(line + "\n")
			if strings.Contains(line, "*/") {
				break
			}
		} else if !foundPrompt {
			output.WriteString(line + "\n")
		}
	}

	if foundPrompt {
		err = ioutil.WriteFile(filename, output.Bytes(), 0644)
		if err != nil {
			fmt.Printf("Error writing to file %s: %v\n", filename, err)
		}
	} else {
		fmt.Printf("No 'Prompt:' comment block found in %s.\n", filename)
	}
}