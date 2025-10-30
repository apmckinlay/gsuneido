### Assert
<pre>
(condition, msg = "")
(expression <i>matcher</i>: value, msg = "")
</pre>

If the condition is true, Assert does nothing. If the condition is false, Assert throws an exception of "Assert FAILED: " $ description.

The second form of Assert uses "matchers" to make the error messages more readable.

For example:

``` suneido
Assert(x is 123)
```

can be written as:

``` suneido
Assert(x is: 123)
```

with the advantage that if it fails the error will be: "expected 123 but it was ..." instead of just "assertion failed"

The msg argument can be used to add a message to the exception to help identify where it came from.

**Note:** With new style matcher Assert's the msg argument must be named.

``` suneido
Assert(x is: 123, msg: "count to three")
```

If the expression is a block, then it will be executed with Catch. This is useful for catching exceptions:

``` suneido
Assert({ Object().x } throws: "member not found")
```

The matchers include:
`is: value`
: 

`isnt: value`
: 

`matches: value`
: Regular expression match.

`has: value`
: Checks if a string or object .Has?(value) returns true

`hasnt: value`
: Checks if a string or object .Has?(value) returns false.

`hasMember: value`
: Checks if an object .Member?(value) returns true

`hasntMember: value`
: Checks if an object .Member?(value) returns false

`hasAssoc: #(member: value)`
: Checks if an object contains the specified association

`like: string`
: Trim's leading and trailing whitespace, normalizes line endings by removing '\r', and converts sequences of spaces or tabs to single spaces before comparing.

`startsWith: string`
: 

`endsWith: string`
: 

`lessThan: value`
: 

`lessThanOrEqualTo: value`
: 

`greaterThan: value`
: 

`lessThanOrEqualTo: value`
: 

`greaterThanOrEqualTo: value`
: 

`between: #(a, b)`
: Confirms that the supplied number is between a and b (*i.e.*`a <= x and x <= b`).

`closeTo: #(value, decimalPlaces)`
: Confirms that the supplied number, when rounded to `decimalPlaces` decimal places, is equal to a benchmark `value` also rounded to `decimalPlaces` decimal places.

`greaterThanOrEqualTo: value`
: 

`isNumber: value`
: 

`isString: value`
: 

`isObject: value`
: 

`isType: type`
: Checks that Type(value) is type

`throws: string`
: Confirms that the supplied block throws an exception that contains (.Has?) the given string. For more complex exception checking, use 
[Catch](<Catch.md>).   
e.g. `Assert({ [].Max() } throws: "empty")`

Another advantage to matchers is that you can write your own. This can make the asserts simpler and easier to read. A class named Matcher_... will be automatically recognized as a matcher. Matchers define Match, Expected, and Actual methods. Look at the existing matchers in stdlib for examples. Matchers are also used by [Mock](<Mock.md>)