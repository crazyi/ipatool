# AGENTS.md

This file provides guidance for coding agents working in this repository.

## Repository overview
- Project: `ipatool`
- Language: Go (1.23+)
- Entry point: `main.go`
- CLI command implementations: `cmd/`
- Core business logic: `pkg/appstore/`
- HTTP layer: `pkg/http/`
- Utilities: `pkg/util/`, `pkg/log/`, `pkg/keychain/`

## Build, Lint, and Test Commands

### Local checks
Run these commands from the repository root before submitting changes:

```bash
# Generate mocks (required after interface changes)
go generate ./...

# Run all tests
go test ./...

# Build all packages
go build ./...

# Run linter
golangci-lint run ./...

# Run a single test
go test -v ./pkg/appstore/... -run "TestName"

# Run tests with verbose output
go test -v ./...

# Run tests in watch mode (requires gotestsum or similar)
go test ./... -count=1
```

### Code generation
- Mocks are generated using `go:generate` with `mockgen`
- Run `go generate ./...` whenever interfaces in `pkg/http/client.go` or other files change
- Generated mocks are placed alongside tests (e.g., `client_mock.go`)

## Development workflow
1. Keep changes focused and minimal.
2. Prefer idiomatic Go and keep command behavior consistent with existing commands in `cmd/`.
3. Run formatting, linting, and tests before finalizing changes.
4. Use dependency injection for testability (see `cmd/common.go` for the `Dependencies` pattern).

## Code Style Guidelines

### Formatting
- Use `gofmt` and `goimports` for code formatting (enabled in `.golangci.yml`).
- Run `gofmt -w .` before committing.
- Imports are grouped: standard library first, then third-party packages (blank line between groups).

### Naming Conventions
- **Variables**: camelCase (e.g., `appID`, `bundleID`)
- **Constants**: PascalCase for exported, camelCase for unexported (e.g., `KeychainServiceName`, `configDirectoryName`)
- **Functions**: PascalCase for exported (e.g., `SearchCmd`), camelCase for unexported (e.g., `newLogger`)
- **Types/Interfaces**: PascalCase (e.g., `AppStore`, `SearchInput`, `SearchOutput`)
- **Packages**: short, lowercase, no underscores (e.g., `pkg/http`, `pkg/appstore`)
- **Files**: lowercase with underscores for multiple words (e.g., `appstore_search_test.go`)

### Interfaces
- Define interfaces in the same package that uses them or in a dedicated package.
- Use small, focused interfaces (single responsibility).
- Name interfaces with the suffix `er` or descriptive names (e.g., `AppStore`, `Client`).
- Place interfaces before implementing structs.

### Structs and Types
- Use input/output structs for operations (e.g., `SearchInput`, `SearchOutput`).
- Embed interfaces for dependency injection (e.g., `type appstore struct { keychain keychain.Keychain }`).
- Use tags for serialization (e.g., `json:"resultCount,omitempty"`).

### Error Handling
- Use `fmt.Errorf` with `%w` for wrapping errors: `return nil, fmt.Errorf("failed to do something: %w", err)`.
- Return errors directly from called functions when no additional context is needed.
- Use sentinel errors for known failure cases (e.g., `var ErrPasswordTokenExpired = errors.New(...)`).
- Use `errors.Is()` for error checking in callers.
- Add `// nolint:wrapcheck` comment for intentional error passthrough in CLI commands.

### Logging
- Use zerolog for structured logging.
- Use `.Log().Send()` for main log output.
- Use `.Verbose()` for debug/verbose logging.
- Chain fields using `.Str()`, `.Int()`, `.Bool()`, etc.

### Testing
- Use Ginkgo/Gomega for testing framework.
- Use `go.uber.org/mock` (mockgen) for generating mocks.
- Follow the pattern in `pkg/appstore/appstore_search_test.go`:
  - Use `Describe` and `When` blocks for test organization.
  - Use `BeforeEach` for setup.
  - Use `It` for individual test cases.
  - Use `Expect(err).ToNot(HaveOccurred())` for error assertions.
  - Use `gomock` for mocking interfaces.
- Test file naming: `<package>_<feature>_test.go`

### CLI Commands (cmd/)
- Use Cobra for CLI framework.
- Use `RunE` for commands that return errors.
- Group flags in a `var` block at the beginning of the function.
- Use descriptive flag names: `--app-id`, `--bundle-identifier`, `--output`.
- Return meaningful error messages to users.

### Dependencies
- Avoid introducing new dependencies unless necessary.
- Use existing patterns in the codebase before adding new libraries.
- Tool dependencies are defined in `tools.go`.

### Additional Linter Rules (.golangci.yml)
The project enables these additional linters:
- `ginkgolinter`, `godot`, `godox`, `importas`, `nlreturn`, `nonamedreturns`
- `prealloc`, `predeclared`, `unconvert`, `unparam`, `usestdlibvars`
- `wastedassign`, `wrapcheck`, `wsl`

## Commit/PR Guidance
- Write clear, scoped commit messages.
- Use conventional commit format: `type(scope): description`
- Types: `feat`, `fix`, `docs`, `style`, `refactor`, `test`, `chore`
- Summarize what changed and why in PR descriptions.
- Include test/build results in your handoff.
