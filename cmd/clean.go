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
 1. Find the FIRST comment block in the file that has the string "Prompt:" as its first contents (on the same line or on the following one).
 2. Find the end of that comment block.
 3. Delete anything past that comment block.
 4. Overwrite the original file with the new contents.

Hints:
- Don't bother supporting // style comments.
- There's no such thing as strings.NewScanner
*/