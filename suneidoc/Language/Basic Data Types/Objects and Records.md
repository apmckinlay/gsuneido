### Objects and Records

Objects are Suneido's general-purpose containers. They can be used as arrays, lists, records. Objects are also used to represent instances of classes.

Objects only support `is`, `isnt`, `<`, `<=`, `>`, `>=`, subscript `[]`, and member (.) operations.

Objects are equal if they have the same values and members. Comparisons `<`, `<=`, `>`, `>=`

Additional, user defined methods can be added in a class called "Objects" (or "Records").

Object constants can be written as follows:

object-constant:
<pre>   # ( <i>members</i> )</pre>

list-member:

``` suneido
    constant
```

named-members:

``` suneido
    name: constant
```

name:

``` suneido
    identifier
    number
    string
```

For example:

``` suneido
#(1, 2, "abc", "def")
#(name: "Joe", age: 23)
#(1, "abc", name: "Joe")
```

Nested object constants do not require the leading '#' for example:

``` suneido
#(name: "Joe", children: ("Sue", "Sam"))
```

If a string is a valid identifier the quotes are optional. The following is equivalent to the last example:

``` suneido
#(name: Joe, children: (Sue, Sam))
```

String concatenation ($) is allowed within object constants. This concatenation is done at compile time.

``` suneido
#(message: "now is the time\n" $
    "for all good men\n")
```

Record constants are written similarly, except with curly braces instead of parenthesis:

``` suneido
#{name: "Joe", age: 23}
```

Objects and Records can be created at run-time with Object(...) or Record(...) or with [...]

[...] when it has **unnamed** members, is a shortcut for Object(...)

``` suneido
Type([1, 2, 3]) => "Object"
Type([1, 2, a: 3]) => "Object"
```

[...] with **only** named members (or no members) is a shortcut for Record(...)

``` suneido
Type([a: 1, b: 2]) => "Record"
```

See also:
[Object](<../Reference/Object.md>),
[Record](<../../Database/Reference/Record.md>)