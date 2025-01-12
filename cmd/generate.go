package cmd

/*
Prompt:
- Use github.com/spf13/cobra
- Create an external command var called GenerateCmd.
- Full name of the command is 'generate', but it also is available as an alias 'gen'.
- By default this command takes one or more filenames. This is indicated when running --help as [FILES].
- Alternatively it takes one argument: --all (alias -a), in which case no filenames are allowed.
  In this case it scans the current directory and its subdirectories for files matching the "*.go" filename pattern.
- If no arguments are provided, print help for the argument.
- With the available file paths (whether provided or discovered), the following logic is executed:
  pre, filePrompt, post, err = utils.GoParseFile(path)
  prompt = promptPrefix + filePrompt
  generateFile(path, pre, prompt, post)

Specifications:
generateFile(path, pre, prompt, post) executes the following logic:
  content = prompt + "\n\nCODE BLOCK START\n\n" + pre + post + "\n\nCODE BLOCK END"
  use `content` as `prompt` for providers.QueryOpenAI()
  If first line of response starts with ```, strip it.
  If last line starts with ```, strip it as well.
  Overwrite file given in `path` with prompt + '\n' + response.
*/

import (
	"aimake/providers"
	"aimake/utils"
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

const promptPrefix = `
I need to write computer code. Do not add any explanations before or after the code. Output only code.

General requirements for me code are:
Write code that is complete and directly runnable.
DO NOT omit code or use comments such as "more content here" or "code remains unchanged."
NEVER change anything above the line that contains the string "END IMMUTABLE CODE BLOCK"
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
					prompt := promptPrefix + filePrompt
					generateFile(path, pre, prompt, post)
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
				prompt := promptPrefix + filePrompt
				generateFile(filename, pre, prompt, post)
			}
		}
	},
}

func init() {
	GenerateCmd.Flags().BoolP("all", "a", false, "Generate all .go files in the current directory and subdirectories")
}

func generateFile(path, pre, prompt, post string) {
	content := prompt + "\n\nCODE BLOCK START\n\n" + pre + post + "\n\nCODE BLOCK END"
	response, err := providers.QueryOpenAI("gpt-4o-mini", content, "", "")
	if err != nil {
		fmt.Println("Error querying OpenAI:", err)
		return
	}
	fmt.Println(response)
}
