# Contributing

Thanks for helping improve Visum! This project prioritizes clarity, testability, and Go idioms.

## Workflow
1. Create a focused branch for each change.
2. Keep commits small and scoped.
3. Update docs and tests whenever behavior changes.

## Code Style
- Follow `gofmt` and standard Go conventions.
- Keep `internal/core` pure and deterministic.
- Avoid DOM or WebAssembly dependencies outside `internal/adapter/web`.
- Prefer explicit state transitions over hidden side effects.

## Testing
Run tests before submitting:

```
go test ./internal/...
```

## Documentation
- Update `README.md` if any user-facing behavior changes.
- Add comments for exported types and functions.

## Issue Reports
Include:
- Steps to reproduce
- Expected vs. actual behavior
- Browser and Go version

## Pull Requests
- Describe the motivation and scope.
- Link related issues if relevant.
- Note any follow-up tasks or limitations.
