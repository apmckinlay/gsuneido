## Regular Expressions

Regular expressions are used for the =~ and !~ 
[comparison operators](<Expressions/Comparison.md>),
and also for
[string.Match](<Reference/String/string.Match.md>),
[string.Extract](<Reference/String/string.Extract.md>) and
[string.Replace](<Reference/String/string.Replace.md>).

### Matching Ordinary Text

The simplest aspect of *regular expressions* is matching ordinary text.  For example:

``` suneido
"now is the time" =~ "the" 
    => true
```

Notice that with the `=~` operator, 
the right hand *pattern* string only has to be contained in the left hand text string.
The pattern does not have to match the entire string.

The `!~` operator is similar to  `!=`
giving false when the pattern matches, and true when it doesn't.

### Start and End of Line

`^` (caret) and `$` (dollar) 
allow you to match the beginning and end, respectively,
of a line of text.

For example, `"^The"` will match "The" at the beginning of a line.
Similarly, `"it$"` matches "it" at the end of a line.

Note: The caret and dollar match *positions* rather than actual text characters.

Use both caret at the beginning and dollar at the end
if you want to match an entire line.
For example:

``` suneido
"now is the time" =~ "the time" => true
"now is the time" =~ "^the time$" => false
"the time" =~ "^the time$" => true
```

Suneido matches caret either at the start of a string
or after a newline.  
Similarly, dollar matches either at the end of a string
or before a newline (or return).
If you specifically want to match the beginning or end of the entire string,
regardless of newlines, use \A and \Z.

Note: To turn off the special meaning of a regular expression character,
preceed it with a backslash.
For example, use `\^` to match a literal caret.

### Character Classes

The `[`...`]` construct, called a *character class*,
lets you list several characters that are allowed at that point.
For example, to match `"the"` or `"The"`
you could use `"[tT]he"`.

A character class can contain any number of characters.
For example, `[aeiou]` would match a vowel.

Within a character class, a `'-'` (dash) can be used to specify a range of characters.
For example, `[0-9]` for a digit,
or `[a-z]` for a lower case letter.
You can use more than one range in a character class,
for example, `[a-zA-Z]` for an upper or lower case letter,
or `[0-9a-fA-F]` for a hexadecimal digit.
Ranges can be combined with individual characters
as in `[a-zA-Z_.!?]`.

A dash is only special within a character class 
and not the first or last character in a character class.
For example, `[-z]` matches a literal dash or a 'z'.

A *negated* character class matches any character that is *not* listed.
It is written like `[^`...`]`.
For example, to find a 'q' followed by any character *except* a 'u' you could use
`"q[^u]"`.

The only characters that are special within a character class are:

|  |  | 
| :---- | :---- |
| `^` | at the beginning to specify a negated character class | 
| `-` | between two characters to specify a range | 
| `]` | to mark the end of the character class | 
| `\w \W \d \D \s \S` | character class shortcuts (see below) | 
| `[:...:]` | posix character class shortcuts (see below) | 


To include a literal ']' in a character class, make it the first character,
for example `[][]` will match an opening or closing square bracket,
and `[^][]` will match any character 
*except* an opening or closing square bracket.

Note: A character class containing any single character 
(except for '^')
matches that character literally.
This means you can "turn off" the special meaning of a character
by enclosing it in square brackets.

### Character Class Shortcuts

The following shortcuts can be used in place of their equivalent character classes:

|  |  |  | 
| :---- | :---- | :---- |
| `\d` | `[0-9]` | digit | 
| `\D` | `[^0-9]` | non-digit | 
| `\s` | `[\x09-\x0d\x20]` | whitespace | 
| `\S` | `[^\x09-\x0d\x20]` | non-whitespace | 
| `\w` | `[a-zA-Z0-9_]` | part of word | 
| `\W` | `[^a-zA-Z0-9_]` | non-word character | 


For example, you could use `"\d\d\d\d\d"` to match a five digit zipcode
(rather than `"[0-9][0-9][0-9][0-9][0-9]"`).

`\s` includes linefeed and return.
(And `\S` excludes them.)

### Posix Character Classes

|  |  | 
| :---- | :---- |
| `[:alnum:]` | letters and digits | 
| `[:alpha:]` | letters (a-zA-Z) | 
| `[:blank:]` | space or tab only | 
| `[:cntrl:]` | control characters | 
| `[:digit:]` | decimal digits (0-9) | 
| `[:graph:]` | printing characters, excluding space | 
| `[:lower:]` | lower case letters (a-z) | 
| `[:print:]` | printing characters, including space | 
| `[:punct:]` | printing characters, excluding letters (a-zA-Z) and digits (0-9) | 
| `[:space:]` | whitespace | 
| `[:upper:]` | uppercase letters (A-Z) | 
| `[:xdigit:]` | hexadecimal digits (0-9a-fA-F) | 


Unlike the character class shortcuts, posix character classes can only be used **within** a character class. For example:

``` suneido
"a2" =~ "[[:digit:]]"
    => true
```

### Matching Any Character - Dot

A '.' or *dot* will match any single character, 
except for a newline, return, or nul.

For example, to match dates like "01/02/03" or "01-02-03",
you could use `"\d\d.\d\d.\d\d"`.
(Although in practice this might be too general; for example, it would match "12345678".)

Note: dot is equivalent to `[^\r\n\x00]`.

### Matching One of Several Choices - Alternation

The '`|`' character allows you to match one of several alternatives.
For example, we could match either "first" or "second" using
`"first|second"`.
To *constrain* the alternation, use parenthesis, as in
`"(first|second) place"`,
which would match either "first place" or "second place".

Be careful when using alternation with caret and dollars
(and other similar matching).
For example, in `"^From|To"` the alternatives are
`"^From"` and `"To"`.
Usually what you want is `"^(From|To)"`
so the caret is not just part of the first alternative

Note: Alternation will match the *first* alternative
that allows the entire match to succeed.
This is not necessarily the longest match. 
For example:

``` suneido
"the category is basic".Extract("cat|category")  =>  "cat"
"the category is basic".Extract("category|cat")  =>  "category"
```

### Optional and Repeated Items - Quantifiers

The '`?`' character means the preceding item is *optional*.
For example, in <code>"Jul<b>y?</b>"</code> the 'y' is optional,
so it will match either "Jul" or "July".

To make multiple characters optional you can apply the question mark to 
a parenthesized group.  
For example: `"1(st)?"` will match either "1" or "1st".

Similar to the question mark,
'`+`' (*plus*) means "one or more of the preceding item" and
'`*`' (*star*) means "optional or one or more of the preceding item".
Both plus and star try to match as many times as possible,
while still allowing anything afterwards to match.
Like the question mark, plus and star apply to either
the preceding character, or the preceding parenthesized group.

For example, `"hello.*world"` will match  "hello" followed by "world"
with anything (or nothing) in between.
(Since dot doesn't match newlines, the "hello" and the "world" would have to be on the same line.)

Or, `"[0-9]+"` will match one or more digits.

When using patterns like `".*"` remember that 
it will match as much as possible.  For example:

``` suneido
"<one> <two>".Extract("<.*>")  =>  "<one> <two>"
```

One way around this is to use something like:

``` suneido
"<one> <two>".Extract("<[^>]*>")  =>  "<one>"
```

However, this won't work for more than a single character.
A better solution is to use the "non-greedy" versions.
Normally, `?`, `*`, and `+` are "greedy",
meaning they match as much as possible.  
Following them with a question mark, 
`??`, `*?`, and `+?`
makes them "non-greedy", 
meaning they match as little as possible.
For example:

``` suneido
"<one> <two>".Extract("<.*?>")  =>  "<one>"
```

Since dot (`.`) does not match newlines,
if you want your match to span multiple lines you have to use something like
`[^\000]*` where `[^\000]` means any character except NUL (0).

### Ignore Case

You can specify that case should be ignored ('a' will match 'A' and vice versa) by the sequence `(?i)` (ignore case). You can turn this off with `(?-i)`.

Note: As of 2014-12-18 ignore case only applies to a-zA-Z (not extended ascii).

### Word Boundaries

Just like `^` and `$` match the beginning and end of lines, 
`\<` and `\>` match the beginning and end of *words*.
For example: `"\<cat\>"` would match "cat"
but not "scat" or "catalog".

"Start of word" and "end of word" are simply positions 
where a sequence of "word" characters (alphanumeric or underline `[a-zA-Z0-9_]`) 
begins or ends.

### Backreferences

not supported

### Escapes for Special Characters

When Suneido compiles string literals, it converts *escape* sequences.

|  |  | 
| :---- | :---- |
| `\t` | tab | 
| `\n` | newline (linefeed) | 
| `\r` | carriage return | 
| `\0` | nul (zero) character | 
| `\xhh` | an 'x' followed by a two digit hexadecimal number e.g. `\x0a` or `\x0A` | 
| `\"` | a literal double quote | 


### Literal Matching

The special sequence `(?q)` will turn off the meaning of the special characters described above. The only special sequence recognized is `(?-q)` which turns on the meaning of special characters.