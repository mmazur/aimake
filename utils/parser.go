/*
Package "utils".

Functions:
ParseGoFile takes one parameter `path`, which is a path to a golang file.
It parses the file and returns:
- `filedoc` the file's godoc comment
- `code` the code in the file past the first comment block
- `error` an optional error value.
Its logic is:
1. Read the file specified in `path`.
2. Check if the file starts with a comment (comments starting with /* and // must be supported).
- If yes, parse the file until you find the end of the comment block, then set the `filedoc` value to the comment block.
- If not, leave `filedoc` empty.
- `filedoc` is to contain the contents of the first comment block at the start of the file. Never add any comments to it that are further down in the file.
3. Put everything after the `filedoc` comment block into `code`.
4. Before returning normalize `filedoc`:
- If it's a // comment block, those two characters are stripped from each line.
- If it's a /* comment block, both it and its comment-block-ending equivalent are stripped.
5. Before returning normalize `code` by stripping any whitespace from the start and end of it.
*/

package utils

import (
	"bufio"
	"os"
	"strings"
)

// ParseGoFile parses a Go file and extracts the godoc comment and code.
func ParseGoFile(path string) (filedoc string, code string, err error) {
	file, err := os.Open(path)
	if err != nil {
		return "", "", err
	}
	defer file.Close()

	var (
		commentBuilder strings.Builder
		codeBuilder    strings.Builder
		foundDoc       bool
		inBlockComment bool
	)

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		trimmed := strings.TrimSpace(line)

		if !foundDoc {
			// Check for block comment start
			if strings.HasPrefix(trimmed, "/*") {
				inBlockComment = true
				foundDoc = true

				// Add the text after "/*" to the comment builder
				after := strings.TrimPrefix(trimmed, "/*")
				after = strings.TrimSpace(after)
				if strings.Contains(after, "*/") {
					parts := strings.SplitN(after, "*/", 2)
					commentBuilder.WriteString(parts[0] + "\n")
					inBlockComment = false
					// Everything after "*/" is code
					codeBuilder.WriteString(parts[1] + "\n")
				} else {
					commentBuilder.WriteString(after + "\n")
				}
				continue
			}

			// Check for single-line comment
			if strings.HasPrefix(trimmed, "//") {
				streamlined := strings.TrimPrefix(trimmed, "//")
				commentBuilder.WriteString(streamlined + "\n")
				foundDoc = true
				continue
			}

			// If this line is empty, keep scanning for a possible top comment
			if trimmed == "" {
				continue
			}

			// If we reach here, we found a non-comment line => no doc at the top
			// from now on, everything is code
			foundDoc = true
			codeBuilder.WriteString(line + "\n")
			continue
		}

		// If we are actively parsing a block comment
		if inBlockComment {
			if strings.Contains(trimmed, "*/") {
				inBlockComment = false
				parts := strings.SplitN(trimmed, "*/", 2)
				commentBuilder.WriteString(parts[0] + "\n")
				// The rest of the line goes to code
				codeBuilder.WriteString(parts[1] + "\n")
			} else {
				commentBuilder.WriteString(line + "\n")
			}
			continue
		}

		// If doc is found or concluded, everything else is code
		codeBuilder.WriteString(line + "\n")
	}

	if err := scanner.Err(); err != nil {
		return "", "", err
	}

	filedoc = normalizeComment(commentBuilder.String())
	code = strings.TrimSpace(codeBuilder.String())
	return filedoc, code, nil
}

// normalizeComment cleans up the first doc comment by removing /*, */, //, and leading/trailing whitespace.
func normalizeComment(rawComment string) string {
	// Remove any trailing "*/" if it exists (e.g. if the block ended exactly on the line).
	rawComment = strings.ReplaceAll(rawComment, "*/", "")
	lines := strings.Split(rawComment, "\n")

	var cleaned []string
	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		// Remove any leading "//"
		if strings.HasPrefix(trimmed, "//") {
			trimmed = strings.TrimSpace(strings.TrimPrefix(trimmed, "//"))
		}
		cleaned = append(cleaned, trimmed)
	}

	result := strings.Join(cleaned, "\n")
	return strings.TrimSpace(result)
}
