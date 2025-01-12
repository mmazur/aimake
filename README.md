# aimake
Make AI make code  
… while you just tell it where and how.

## This code is terrible

Yes.

## Usage

### Making edits

1. Edit prompt in some_file.
2. `aimake generate some_file`

### Regenerating the whole codebase

1. `aimake clean --all`
2. Not sure yet.

# Development

## Contributing

- Until we get to Milestone 2, it won't be possible to fully develop aimake using aimake. What I'm doing instead is using the prompt in `copilot.txt` together with GitHub Copilot to simulate the `generate` feature.
- PRs must contain prompts that generated the code. Human-only code edits would be cheating.
  - You can pretend code changes have been fully generated using the attached prompt. I won't notice.

## License

Do licenses make sense if the whole thing is AI-generated? And AI-regeneratable?

What happens if you AI-rewrite the existing prompts and then AI-generate fresh code based on them?

## Milestones

### M1 The Basics
- [x] Working `clean`
- [ ] Working `generate`, at least for basic use cases

### M2 Next Generation
- [ ] `generate` is advanced enough to regenerate all files in this project, including ones with external dependencies (like `main.go`)

## Feature Ideas
- [ ] Multi-shot `gen` (-i arg?); feeds the prompt+code N times with same or slightly modified prompts
  - Defaults to >=2 for more reliable `gen`?
- [ ] Support other models (anthropic, self-hosted, etc.)
- [ ] Come up with automatic evaluation of generated code
- [ ] Once automatic evaluation works, come up with automatic iterative optimization (generate code in multiple ways, automatically figure out which seems like the best one) 
- [ ] Automated multi-shot if first try doesn't build/work (related to previous point)
- [ ] Automatic dependency analysis and interface extraction for `generate`
  - Example: `generate main.go` notices it depends on a function in `cmd.go`, that function declaration is extracted and injected into the prompt for generating `main.go`
- [ ] Full codebase analysis for `generate --all` that can figure out the file generation order (starts with standalone files with no external dependencies)
- [ ] Single-shot prompts: `generate --prompt` or maybe a separate command like `go`
