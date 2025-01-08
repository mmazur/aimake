package cmd

/*
Prompt:
- Use github.com/spf13/cobra
- Create an external command var called QueryCmd.
- This command runs the providers.QueryOpenAI() function.
- It takes the same parameters as that function, but makes sure the optional parameters are optional.
*/

import (
	"aimake/providers"
	"fmt"

	"github.com/spf13/cobra"
)

// QueryCmd represents the query command
var QueryCmd = &cobra.Command{
	Use:   "query",
	Short: "Execute the QueryOpenAI function",
	Run: func(cmd *cobra.Command, args []string) {
		model, _ := cmd.Flags().GetString("model")
		query, _ := cmd.Flags().GetString("query")
		apiKey, _ := cmd.Flags().GetString("apiKey")
		systemPrompt, _ := cmd.Flags().GetString("systemPrompt")

		response, err := providers.QueryOpenAI(model, query, apiKey, systemPrompt)
		if err != nil {
			fmt.Println("Error:", err)
			return
		}
		fmt.Println(response)
	},
}

func init() {
	QueryCmd.Flags().StringP("model", "m", "gpt-4o-mini", "Name of the LLM model")
	QueryCmd.Flags().StringP("query", "q", "", "The query to send to OpenAI")
	QueryCmd.Flags().StringP("apiKey", "k", "", "OpenAI API key")
	QueryCmd.Flags().StringP("systemPrompt", "s", "", "Optional system prompt")
	QueryCmd.MarkFlagRequired("query")
}
