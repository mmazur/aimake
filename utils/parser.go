package utils

/*
Prompt:
- Lets specify a function GoParseFile() that takes one parameter `path` which is a golang file. Its logic is:
  1. Reads the file specified in `path`.
  2. Looks for the first comment block that starts with the "Prompt:" string.
  3. Returns three strings: pre, prompt, post.
     - `pre` is everything before the comment block
	 - `prompt` is the contents of the comment block
	 - `post` is everything after the comment block
  4. Before returning from the function, make sure that:
     - `pre` doesn't end with the comment start marker
	 - `prompt` doesn't contain any comment start or end markers
	 - `post` doesn't begin with a comment start marker
*/

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func GoParseFile(path string) (string, string, string, error) {
	file, err := os.Open(path)
	if err != nil {
		return "", "", "", fmt.Errorf("error opening file: %v", err)
	}
	defer file.Close()

	var preLines, promptLines, postLines []string
	scanner := bufio.NewScanner(file)
	inPromptBlock := false
	promptFound := false

	for scanner.Scan() {
		line := scanner.Text()
		if strings.Contains(line, "Prompt:") && !promptFound {
			inPromptBlock = true
			promptFound = true
		}

		if inPromptBlock {
			promptLines = append(promptLines, line)
			if strings.TrimSpace(line) == "*/" {
				inPromptBlock = false
			}
		} else if promptFound {
			postLines = append(postLines, line)
		} else {
			preLines = append(preLines, line)
		}
	}

	if err := scanner.Err(); err != nil {
		return "", "", "", fmt.Errorf("error reading file: %v", err)
	}

	pre := strings.Join(preLines, "\n")
	if strings.HasSuffix(pre, "/*") {
		pre = strings.TrimSuffix(pre, "/*")
	}

	if len(promptLines) > 0 {
		promptLines = promptLines[1 : len(promptLines)-1]
	}
	prompt := strings.Join(promptLines, "\n")
	prompt = strings.ReplaceAll(prompt, "/*", "")
	prompt = strings.ReplaceAll(prompt, "*/", "")
	post := strings.Join(postLines, "\n")
	if strings.HasPrefix(post, "*/") {
		post = strings.TrimPrefix(post, "*/")
	}

	return pre, prompt, post, nil
}
