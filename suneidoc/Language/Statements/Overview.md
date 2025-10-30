### Overview

Suneido includes the same basic statements as C, C++, or Java, with some minor additions (forever) and some minor differences (switch).

Wherever a statement is expected, you can use a *compound statement* i.e. zero or more statements enclosed in curly braces.
<pre>
{
<i>statements</i>
}
</pre>

Unlike C or C++, semicolons are only required in order to put multiple statements on a line. For example:

``` suneido
x = 5 ; y = 6
```

Unlike C, C++, or Java, parenthesis are not required around the expressions of control statements if they are terminated with a newline. For example:

``` suneido
if x is false
    return false
```

However, parenthesis are required around a C style "for":

``` suneido
for (i = 0; i < 10; ++i)
```

Omitting parenthesis can lead to ambiguities - see [Whitespace](<../Whitespace.md>).