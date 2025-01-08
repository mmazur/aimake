# aimake
Make AI make your code

## Usage

### Making edits

1. Edit prompt in some_file.
2. `aimake generate some_file`

### Regenerating the whole codebase

1. `aimake clean --all`
2. Not sure yet.

## Contributing

PRs must contain prompts that generated the code. Manual edits would be cheating.

## License

Do licenses make sense if the whole thing is AI-generated? And AI-regeneratable?

## Milestones

### M1 Regenerate Yourself
- [x] Working `clean`
- [ ] Working `generate`

### Other ideas
- [ ] Automatic dependency analysis for `generate --all`