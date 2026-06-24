---
name: suneido-language
description: Reference for Suneido language syntax and semantics
---

# Summary of Suneido Language

For more detailed information see the documentation in the suneidoc folder.
The standard library source code is in the stdlib folder.

Uncapitalized name are local variables.
Capitalized names are global and are only defined in library tables like stdlib.
Global name values are constants.

Expressions:
- use `is, isnt, and, or, not` instead of `==, !=, &&, ||, !`
- regular expression matching is done with `=~` and `!~`
- supports `?:` ternary operator
- use .Size() to get the size of a string or object

"Object" is the general purpose container
- unnamed members start at 0 e.g. `#(a,b,c)[0] is 'a'`
- it can have named and/or unnamed members e.g. `Object(1, 2, a: 3, b: 4)`
- #(...) is an immutable object constant
- Object(...) or Record(...) are mutable
- use ob.Add(...) to append to the unnamed members

anonymous function: `function(parameters) { ... }`

closure: `{|parameters| ... }`
- `return` returns from the containing function, not just the closure
- closures return the value of their last statement
- closures are also known as "blocks"
- break and continue from a closure throw "block:break" and "block:continue"

function and closure parameters:
- can have defaults e.g. `function(a, b, c = 1, d = 2)`
- can use `@var` to receive all the arguments into an object

call arguments:
- can have named arguments e.g. `fn(a, b, c: 3, d: 4)`
- `:var` is shorthand for `var: var`
- extra named arguments for which there are no matching parameters are ignored
- can use `@object` to pass the contents of object as separate arguments

control statements: if-else, while, do-while, for, switch
- switch statements do not fallthrough so do not require break
- if there is no default and none of the cases match it will throw

try-catch: try <statement> catch (var, prefix) statement
- catch is optional
- var and prefix are optional

classes:
```
class : Base
    {
    X: 123
    Add(x) { return x + .X}
    }
```
- classes are immutable
- `: Base` is optional
- capitalized members are public, uncapitalized are private
- `class : Base` can be abbreviated to just `Base`
- classes can have a `CallClass` method
- there is a default `CallClass` that creates an instance
- classes can have a `Call` method to make instances callable
- classes can have a `New` method constructor
- `New` methods can call `super` to specify the arguments
- within class methods, `.name` is shorthand for `this.name`
