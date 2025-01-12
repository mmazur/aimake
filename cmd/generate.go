package cmd

/*
Prompt:
- Use github.com/spf13/cobra
- Create an external command var called GenerateCmd.
- Full name of the command is 'generate', but it also is available as an alias 'gen'.
- By default this command takes one or more filenames. This is indicated when running --help as [FILES].
- Alternatively it takes one argument: --all (alias -a), in which case no filenames are allowed.
  In this case it scans the current directory and its subdirectories for files matching the "*.go" filename pattern.
- If no arguments are provided, print help for the command.
- With the available file paths (whether provided or discovered), the following logic is executed:
  pre, filePrompt, post, err = utils.GoParseFile(path)
  generateFile(path, pre, filePrompt, post)

Function specifications:
generateFile(path, pre, prompt, post) executes the following logic:
  content = promptPrefix + prompt + "\n\nCODE BLOCK START\n" + pre + post + "\nCODE BLOCK END"
  use `content` as `prompt` for providers.QueryOpenAI()
  If first line of response starts with ```, strip it.
  If last line starts with ```, strip it as well.
  Overwrite file given in `path` with: "/*\nPrompt:\n" + prompt + "\n*" + "/\n" + response.
*/

import (
	"aimake/providers"
	"aimake/utils"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
)

const promptPrefix = `
I need to write computer code. Do not add any explanations before or after the code. Output only code.

General requirements for my code are:
 Write code that is complete and directly runnable.
 DO NOT omit code or use comments such as "more content here" or "code remains unchanged."
 NEVER change anything above the line that contains the string "END IMMUTABLE CODE BLOCK" if such a line exists.
 Do not unnecessarily remove any comments or code.

The code itself must do the following:
`

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
					pre, filePrompt, post, err := utils.GoParseFile(path)
					if err != nil {
						fmt.Println("Error parsing file:", err)
						return nil
					}
					generateFile(path, pre, filePrompt, post)
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
				pre, filePrompt, post, err := utils.GoParseFile(filename)
				if err != nil {
					fmt.Println("Error parsing file:", err)
					continue
				}
				generateFile(filename, pre, filePrompt, post)
			}
		}
	},
}

func init() {
	GenerateCmd.Flags().BoolP("all", "a", false, "Generate all .go files in the current directory and subdirectories")
}

func generateFile(path, pre, prompt, post string) {
	content := promptPrefix + prompt + "\n\nCODE BLOCK START\n" + pre + post + "\nCODE BLOCK END"
	response, err := providers.QueryOpenAI("o1-mini", content, "", "")
	if err != nil {
		fmt.Println("Error querying OpenAI:", err)
		return
	}

	// Strip leading and trailing ``` from response
	responseLines := strings.Split(response, "\n")
	if len(responseLines) > 0 && strings.HasPrefix(responseLines[0], "```") {
		responseLines = responseLines[1:]
	}
	if len(responseLines) > 0 && strings.HasPrefix(responseLines[len(responseLines)-1], "```") {
		responseLines = responseLines[:len(responseLines)-1]
	}
	response = strings.Join(responseLines, "\n")

	// Overwrite the file with the new content
	newContent := "/*\nPrompt:\n" + prompt + "\n*" + "/\n" + response
	err = os.WriteFile(path, []byte(newContent), 0644)
	if err != nil {
		fmt.Println("Error writing file:", err)
	}
}
