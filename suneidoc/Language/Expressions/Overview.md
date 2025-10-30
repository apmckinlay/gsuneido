### Overview

Suneido includes the same set of operators as C, C++, and Java, with the addition of regular expression matching.

Expressions are composed of values, operators, and [function calls](<Function Calls.md>).

Values are either names or literals (e.g. boolean, numbers, strings, dates, objects).

The order of evaluation (precedence and associativity) of expressions is similar to C or C++.

``` suneido
++ -- not ! + - ~ (unary)
* / %
+ - $
<< >>
is == isnt != =~ !~
< <= > >=
&
^
|
and &&
or ||
=  +=  -=  $=  >>=  <<=  &=  |=  ^=  *=  /=  %= 
?:
|>
```

Unary operators and assignment operators are right-associative (evaluated right to left) all others are left-associative (evaluated from left to right).

For example:

``` suneido
this:               is the same as: 

2 + 3 * 5           2 + (3 * 5) 
1 > 2 and 3 > 4     (1 > 2) and (3 > 4) 
a = b = c           a = (b = c) 
```

Newlines mark the end of an expression statement, unless it is obviously incomplete e.g. inside (), [], {}, or after a binary operator.

**Note**: Unlike C or C++, comma is **not** an operator. It can still be used to have multiple expressions in a for initialization or increment.

**Note**: Constant expressions are evaluated at compile time. This differs slightly from runtime evaluation. For example, `x * 0` will be compiled to 0, whereas runtime evaluation could throw an exception if x is not a number. This also applies to `and, or, &, |`