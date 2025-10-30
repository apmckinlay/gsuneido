## Whitespace

Whitespace (blanks, tabs, newlines) is ignored except as it serves to separate ambiguous sequences of characters.

Newlines mark the end of an expression statement, unless it is obviously incomplete
e.g. inside (), [], {}, or after a binary operator.

A newline is also significant in the following situations:

<div class='table-style table-half-width'>

| This: | Will be interpreted as: | 
| :---- | :---- |
| `if Func(...) { ... }` | `if Func(..., { ... })` | 
| `if Func(...)`<br />`{ ... }` | `if (Func(...))`<br />	`{ ... }` | 
| `if Name { ... }` | `if (Name { ... })` | 
| `if Name` <br />    `{ ... }` | `if (Name)`<br />    `{ ... }` |

</div>