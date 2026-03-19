# Tool Summarize Output Examples

This document shows the output of summarize functions for all tools.

## suneido_read_book

Read from a Suneido book (documentation) table. Returns a JSON object containing the page 'text' and a 'children' array of sub-topic names.

### Parameters

- `book` (required): Name of the book table (e.g. 'suneidoc')
- `path`: The path to the book page. If sub-topics are returned in 'children', append them to this path to dive deeper. (e.g. 'Database/Reference/Query'). Empty or omitted for root.

### Summarize Output Examples

**Args:** `{book: "suneidoc"}`

=> **Read Book** `suneidoc`

**Args:** `{book: "suneidoc", path: "/Database"}`

=> **Read Book** `suneidoc` path:`/Database`

**Args:** `{book: "suneidoc", path: "/Database/Reference/Query"}`

=> **Read Book** `suneidoc` path:`/Database/Reference/Query`

---

## suneido_read_code

Get the source code from a library for a specific name

### Parameters

- `library` (required): Name of the library (e.g. 'stdlib')
- `name` (required): Name of the definition (e.g. 'Alert')
- `start_line`: 1-based line number to start from (default 1)
- `plain`: If true, don't add line numbers (default false)

### Summarize Output Examples

**Args:** `{library: "stdlib", name: "Alert"}`

=> **Read Code** `stdlib` `Alert`

**Args:** `{library: "stdlib", name: "Alert", start_line: 1}`

=> **Read Code** `stdlib` `Alert`

**Args:** `{library: "stdlib", name: "Alert", start_line: 10, plain: true}`

=> **Read Code** `stdlib` `Alert` start-line:`10` plain

---

## suneido_create_code

Create a new library definition. The definition must be valid Suneido code. Returns an error if the definition already exists.

### Parameters

- `library` (required): Name of the library (e.g. 'stdlib')
- `path` (required): Folder path within the library (e.g. 'Debugging/Tests', empty string for root)
- `name` (required): Name of the definition (e.g. 'Alert')
- `code` (required): The source code for the definition

### Summarize Output Examples

**Args:** `{library: "stdlib", path: "", name: "TestFunc", code: "TestFunc()↩	{↩	123↩	…"}`

=> **Create Code** `stdlib` `TestFunc`
```suneido
TestFunc()
	{
	123
	}
```

**Args:** `{path: "Debugging/Tests", name: "TestFunc", code: "TestFunc()↩	{↩	123↩	…", library: "stdlib"}`

=> **Create Code** `stdlib` `Debugging/Tests` `TestFunc`
```suneido
TestFunc()
	{
	123
	}
```

---

## suneido_delete_code

Delete a library definition by name

### Parameters

- `library` (required): Name of the library (e.g. 'stdlib')
- `name` (required): Name of the definition (e.g. 'Alert')

### Summarize Output Examples

**Args:** `{library: "stdlib", name: "OldFunc"}`

=> **Delete Code** `stdlib` `OldFunc`

---

## suneido_edit_code

Modify a Suneido definition by inserting or replacing lines.
This tool is the preferred way to edit existing code.
- Lines are 1-based (matching the output of suneido_read_code)
- Modes:
  - "insert_before": Insert lines of code before the specified line
  - "insert_after": Insert lines of code after the specified line
  - "replace_lines": Replace 'count' lines of code starting at 'line'
- For deletions with replace_lines: Set 'code' to an empty string
- Always call suneido_read_code before this to ensure line numbers are current
- Do NOT include line numbers in the replacement code, just the code itself


### Parameters

- `library` (required): Name of the library (e.g. 'stdlib')
- `name` (required): Name of the definition (e.g. 'Alert')
- `mode` (required): Operation mode: 'insert_before', 'insert_after', or 'replace_lines'
- `line` (required): Line number (1-based)
- `count`: Number of lines to replace (only for replace_lines mode)
- `code` (required): Replacement code

### Summarize Output Examples

**Args:** `{library: "stdlib", name: "Alert", mode: "insert_before", line: 5, code: "new line"}`

=> **Edit Code** `stdlib` `Alert` `insert-before` line: `5`
```suneido
new line
```

**Args:** `{library: "stdlib", name: "Alert", mode: "insert_after", line: 10, code: "new line"}`

=> **Edit Code** `stdlib` `Alert` `insert-after` line: `10`
```suneido
new line
```

**Args:** `{library: "stdlib", name: "Alert", mode: "replace_lines", line: 5, count: 3, code: "replacement"}`

=> **Edit Code** `stdlib` `Alert` `replace-lines` line: `5 to 7`
```suneido
replacement
```

---

## suneido_execute

Executes Suneido code for its result or side effects.
Use this for calculations, data manipulation, or system commands.
A single returned object will appear as the first result (e.g., [[1,2]])
multiple return values appear as separate elements (e.g., [1,2]).
Errors will include the call stack trace

### Parameters

- `code` (required): Suneido code to execute (as the body of a function)

### Summarize Output Examples

**Args:** `{code: "1 + 2"}`

=> **Execute** `1 + 2`

**Args:** `{code: "Date.Now()"}`

=> **Execute** `Date.Now()`

**Args:** `{code: "result = []↩for (i =…"}`

=> **Execute**
```suneido
result = []
for (i = 0; i < 10; i++)
	result.Add(i)
result
```

---

## suneido_check_code

Checks Suneido code for syntax and compilation errors without executing it. Returns compiler warnings only.

### Parameters

- `code` (required): Suneido code to check (as the body of a function)

### Summarize Output Examples

**Args:** `{code: "function foo() { return 123 }"}`

=> **Check Code**
```suneido
function foo() { return 123 }
```

**Args:** `{code: "function bar(x) { x + 1 }"}`

=> **Check Code**
```suneido
function bar(x) { x + 1 }
```

---

## suneido_code_folders

List folders and code items under a library path

### Parameters

- `library` (required): Name of the library (e.g. 'stdlib')
- `path` (required): Folder path within the library (e.g. 'Debugging/Tests', empty string for root)

### Summarize Output Examples

**Args:** `{library: "stdlib", path: ""}`

=> **Code Folders** `stdlib`

**Args:** `{library: "stdlib", path: "Debugging"}`

=> **Code Folders** `stdlib` `Debugging`

**Args:** `{library: "stdlib", path: "Debugging/Tests"}`

=> **Code Folders** `stdlib` `Debugging/Tests`

---

## suneido_libraries

Get a list of the libraries currently in use in Suneido

### Parameters


### Summarize Output Examples

**Args:** `{}`

=> **Libraries**

---

## suneido_query

Execute a Suneido database query and return the results as Suneido-format text (Value.String) in a simple row/column array format (limit 100)

### Parameters

- `query` (required): Suneido query (e.g. 'tables sort table')

### Summarize Output Examples

**Args:** `{query: "tables sort table"}`

=> **Query** `tables sort table`

**Args:** `{query: "columns↩where table …"}`

=> **Query**
```suneido
columns
where table = 'test'
```

---

## suneido_schema

Get the schema for a Suneido database table, or the definition for a view

### Parameters

- `table` (required): Name of the table or view to get schema for

### Summarize Output Examples

**Args:** `{table: "tables"}`

=> **Schema** `tables`

**Args:** `{table: "views"}`

=> **Schema** `views`

---

## suneido_search_book

Search book pages by regex on path and text

### Parameters

- `book` (required): Name of the book table (e.g. 'suneidoc')
- `path`: Regular expression applied to the full page path (path + name)
- `text`: Regular expression applied to page text (optional if path provided)
- `case_sensitive`: If true, regex matching is case sensitive (default false)

### Summarize Output Examples

**Args:** `{book: "suneidoc", path: "Database"}`

=> **Search Book** `suneidoc` path:`Database`

**Args:** `{book: "suneidoc", text: "query"}`

=> **Search Book** `suneidoc` text:`query`

**Args:** `{book: "suneidoc", path: "Database", text: "query", case_sensitive: true}`

=> **Search Book** `suneidoc` path:`Database` text:`query` case-sensitive

---

## suneido_search_code

Search library code by exact library name and regex on name/text

### Parameters

- `library`: Library name (exact match, optional; default all libraries)
- `name`: Regular expression applied to definition names (optional if code provided)
- `code`: Regular expression applied to definition text (optional if name provided)
- `case_sensitive`: If true, regex matching is case sensitive (default false)
- `modified`: If true, only return results where the code has been modified

### Summarize Output Examples

**Args:** `{name: "Alert"}`

=> **Search Code** name:`Alert`

**Args:** `{library: "stdlib", name: "Alert"}`

=> **Search Code** library:`stdlib` name:`Alert`

**Args:** `{library: "stdlib", code: "function"}`

=> **Search Code** library:`stdlib` code:`function`

**Args:** `{library: "", name: "^A", code: "", case_sensitive: true}`

=> **Search Code** name:`^A` case-sensitive

**Args:** `{library: "stdlib", modified: true}`

=> **Search Code** library:`stdlib` modified

---

## suneido_tables

Get a list of database table names that start with the given prefix (limit of 100)

### Parameters

- `prefix` (required): Only return tables whose names start with this prefix (empty string for all)

### Summarize Output Examples

**Args:** `{prefix: ""}`

=> **Tables**

**Args:** `{prefix: "test"}`

=> **Tables** `test`

---

