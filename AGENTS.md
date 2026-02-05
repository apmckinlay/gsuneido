# gSuneido Project Context

gSuneido is the Go implementation of the Suneido integrated language and database.

**Language Reference**: For syntax and semantics, strictly refer to @.kilocode/skills/suneido-language/SKILL.md
**Database Reference**: For query language, strictly refer to @.kilocode/skills/suneido-database/SKILL.md
**Suneido Code**: Stored in library tables (database), compiled to bytecode, and interpreted. Suneido standard library code can be found in `.ss` files under the `stdlib/` directory
**Database**: Immutable append-only with MVCC and relational algebra query language.

## Critical Rules
- **IMPORTANT**: Do not overwrite or delete `gsuneido.exe` or `suneido.db`.
- **Comments**: Minimal comments. Only explain *why*, not *what*.

## Architecture & Key Files
- **Entry Point**: `cmd/gsuneido/` (Main executable)
- **Core** (`core/`):
  - `interp.go`: Main interpreter loop.
  - `value.go`: `Value` interface definition.
  - `builtin/`: Built-in Suneido functions.
  - `compile/`: Parser and code generation.
- **Database**:
  - `db19/`: Append-only database engine (MVCC).
  - `dbms/`: Database server and query engine (`dbms.go`).
- **Suneido Source**: `.ss` files.
- **Data**: `suneido.db` (Database)

## Building
- use `make` to build the project (**not** `go build`)

## Testing
- use `make test` to run all the tests for the project (the makefile already uses `go test -short -timeout=30s`)
- **IMPORTANT** Do not change directories with cd to run tests
- **IMPORTANT** If you run `go test` directly (instead of `make test`), include: `-short -timeout=30s`

## Tooling
- **gopls MCP**: Use the `gopls` MCP server tools to navigate, understand, and verify the Go codebase.
  - Use `go_workspace` to understand the workspace structure.
  - Use `go_search` and `go_symbol_references` for navigation.
  - Use `go_diagnostics` to check for errors after edits.

## Code Style & Conventions
- **Naming**: Suneido values prefixed with `Su` (SuObject, SuStr, SuInt, etc.)
- **License**: MIT license header required
- **Error handling**: Use `panic` for programming errors
- **Types**: Dynamic typing in Suneido, strict typing in Go implementation
- **Tests**: 
  - Use custom `assert.T(t)` helpers e.g. `assert.T(t).This(x).Is(y)`
  - use a test helper function like other tests (**not** data driven)
