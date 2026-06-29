# gSuneido Project Context

gSuneido is the Go implementation of the Suneido integrated language and database.

It uses the latest version of Go (see `/go.mod`)

**Language Skill**: For syntax and semantics, strictly refer to @.agents/skills/suneido-language/SKILL.md (use this skill when you need to look at `/stdlib` `.ss` files )

**Database Skill**: For query language, strictly refer to @.agents/skills/suneido-database/SKILL.md (use this skill when you need to work with database queries in `/dbms/query`)

**Suneido Code**: Stored in library tables (database), compiled to bytecode, and interpreted. Suneido standard library code can be found in `.ss` files under the `stdlib/` directory

**Database**: Immutable append-only with MVCC and relational algebra query language.

## Critical Rules
- **IMPORTANT**: Do not overwrite or delete `gsuneido.exe` or `suneido.db`.
- **Comments**: Concise comments are allowed when they explain *why*, not *what*.

## Architecture & Key Files
- **Entry Point**: `/gsuneido.go`
- **Core** (`core/`):
  - `interp.go`: Main interpreter loop.
  - `value.go`: `Value` interface definition.
  - `builtin/`: Built-in Suneido functions.
  - `compile/`: Parser and code generation.
- **Database**:
  - `db19/`: Append-only database engine (MVCC).
  - `dbms/`: Database server and query engine (`dbms.go`).
  - `dbms/query/`: Query processing and execution.
- **Suneido Source**: `.ss` files under `/stdlib/` 
(just for reference, no easy way to run them)
- **Database**: `suneido.db`

## Building
- use `make` to build the project (**not** `go build`) to get the right build options

## Testing
- **IMPORTANT** Run tests after making changes: go test -short -timeout=30s ./... (changes can affect other packages so it is preferable to run all the tests, not just one package)
- **IMPORTANT** Do not change directories with cd to run tests
- **IMPORTANT** When running `go test` include: `-short -timeout=30s`
- **linting** use `go fmt` and `go vet`

## Tooling
- **gopls MCP**: Use the `gopls` MCP server tools to navigate, understand, and verify the Go codebase.
  - Use `go_workspace` to understand the workspace structure.
  - Use `go_search` and `go_symbol_references` for navigation.
  - Use `go_diagnostics` to check for errors after edits.

## Code Style & Conventions
- **Naming**: Suneido values prefixed with `Su` (SuObject, SuStr, SuInt, etc.)
- **License**: MIT license header required
- **Error handling**: Follow the conventions of the code base for when to use `panic` for programming errors
- **Types**: Dynamic typing in Suneido, strict typing in Go implementation
- **Asserts** for runtime asserts use `assert.That(...)`
- **Tests**: 
  - Use custom `assert.T(t)` helpers, for example `assert.T(t).This(x).Is(y)`
  - use a test helper function like other tests (**not** data driven)
