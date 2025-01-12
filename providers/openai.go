package providers

/*
Prompt:
- This file exposes a single function called `QueryOpenAI` that takes the following text parameters:
	- the name of an llm model; default to "o1-mini" if none provided
	- a query
	- openai api key – if empty, take it from the 'OPENAI_API_KEY' environment variable
	- an optional system prompt – defaults to empty string
- The implementation is to use the github.com/sashabaranov/go-openai library.

Implementation notes:
- Trying github.com/openai/openai-go failed, seems LLMs do not know how to use that library.
*/

import (
	"context"
	"fmt"
	"os"

	openai "github.com/sashabaranov/go-openai"
)

// QueryOpenAI queries the OpenAI API with the given parameters.
func QueryOpenAI(model, query, apiKey, systemPrompt string) (string, error) {
	if apiKey == "" {
		apiKey = os.Getenv("OPENAI_API_KEY")
	}

	if model == "" {
		model = "o1-mini"
	}

	client := openai.NewClient(apiKey)

	messages := []openai.ChatCompletionMessage{
		{
			Role:    openai.ChatMessageRoleUser,
			Content: query,
		},
	}

	if systemPrompt != "" {
		messages = append([]openai.ChatCompletionMessage{
			{
				Role:    openai.ChatMessageRoleSystem,
				Content: systemPrompt,
			},
		}, messages...)
	}

	req := openai.ChatCompletionRequest{
		Model:    model,
		Messages: messages,
	}

	resp, err := client.CreateChatCompletion(context.Background(), req)
	if err != nil {
		return "", fmt.Errorf("failed to query OpenAI: %v", err)
	}

	if len(resp.Choices) == 0 {
		return "", fmt.Errorf("no choices returned from OpenAI")
	}

	return resp.Choices[0].Message.Content, nil
}
