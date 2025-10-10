# Go Style Guide

## Style Guide

- Review and follow the Uber Go Style Guide for all Go code in this repository: [Uber Go Style Guide](https://raw.githubusercontent.com/uber-go/guide/refs/heads/master/style.md)

LLMs working on this project must read this file fully and apply the style guide when writing Go code.

### Additional Guidance for this Repository

- Avoid obvious comments. Do not comment on trivial or idiomatic patterns (e.g., wrapping a context with a timeout, basic error returns, simple getters/setters). Prefer self-explanatory code and comments that capture non-obvious rationale, invariants, or edge cases. This aligns with the spirit of the Uber guide's emphasis on clarity and avoiding clutter. See also: [Uber Go Style Guide](https://raw.githubusercontent.com/uber-go/guide/refs/heads/master/style.md).

## Key Points

- **Error Messages**: Do NOT wrap errors with "failed to" prefixes. Use concise, direct error messages like `"create resource: %w"` instead of `"failed to create resource: %w"`
- **Comments**: Avoid obvious comments. Do not comment on trivial operations like variable assignments, simple error returns, or basic function calls
- **Function Names**: Use clear, descriptive names. Avoid redundant prefixes
- **Code Structure**: Keep functions focused and avoid unnecessary complexity
