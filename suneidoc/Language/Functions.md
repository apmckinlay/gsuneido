## Functions

In Suneido, functions are values. This means they can be assigned to variables, passed to other functions, etc.

Function literals are written as follows:

function: 

``` suneido
function ( parameters ) { statements }
function ( @name ) { statements }
```

parameter:

``` suneido
name
name = literal
```

Parameters must be separated by commas. Parameters with default values must come after parameters without defaults.

``` suneido
f = function (a, b = "") { return a $ b }
f("hello")
    => "hello"
f("hello", " world")
    => "hello world"
```

`@name` places all the supplied arguments into an object. (It must be the only parameter - you cannot specify other parameters as well.) Un-named arguments become un-named members, named arguments become named members. For example:

``` suneido
f = function (@args) { return args }
f(1, 2, a: 3, b: 4)
    => #(1, 2, b: 4, a: 3)
```

A parameter can also be one of the following variations:
`_name`
: [Implicit Dynamic Parameters](<Functions/Implicit Dynamic Parameters.md>)

`.name`
: 

`.Name`
: [Member Parameters](<Classes/Member Parameters.md>)

`._name`
: 

`._Name`
: These are a combination of Dynamic Implicit Parameters plus the shortcut to set members.

Note: This documentation refers to *parameters* when talking about a function definition and *arguments* when talking about the values passed to a function when it is called.

See also: [Function Calls](<Expressions/Function Calls.md>)