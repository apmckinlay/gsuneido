Suneido is an integrated language and database.
gSuneido is the Go implementation of Suneido.

The Suneido language is dynamically typed.
The syntax is similar to C or Java i.e. using curly braces.
Suneido code is stored in library tables in the database.
It is compiled to byte code and the byte code is interpreted.

The database is immutable append-only.
It uses multi-version concurrency control.
It has a relational algebra query language.

Suneido can operate either standalone or client-server.

Do not overwrite gsuneido.exe or suneido.db

Do not add excessive comments.
A comment on almost every line of code is excessive.
Only add comments when they provide information that is not clear from the code.

Changes should follow the style of the existing code.

IMPORTANT: Write tests to reproduce errors BEFORE fixing them.

## Build & Test Commands
- make is preferred over `go build`
- `make port` - Build command-line gsport.exe only
- `make test` - Run Go tests with `-short` flag
- `go test -run TestFunction ./package` - Run specific test function
- `go test ./package` - Run tests in specific package
- `go test -benchmem -bench=BenchmarkName ./package` - run benchmark
- do not change directories to run tests
- `make clean` - Remove built files and clean cache

## Architecture & Structure
- **Core packages**: `core/` (values, types), `builtin/` (built-in functions), `compile/` (parser, codegen), 'util' (miscellaneous utility functions)
- **Database**: `db19/` (append-only DB with MVCC), `dbms/` (query engine, client-server)
- **Suneido code**: `.ss` files contain Suneido source
- **Database files**: `suneido.db` is the actual database, `.su` files are dumped database tables

## Code Style & Conventions
- **Naming**: Suneido values prefixed with `Su` (SuObject, SuStr, SuInt, etc.)
- **License**: MIT license header required
- **Error handling**: Use `panic` for programming errors
- **Types**: Dynamic typing in Suneido, strict typing in Go implementation
- **Tests**: 
  - Use custom `assert.T(t)` helpers e.g. `assert.T(t).This(x).Is(y)`
  - test helper function (not data driven)

Use gsport REPL to run Suneido code.
Remember to "make port" before using it.

When working with Suneido code refer to @suneido.md