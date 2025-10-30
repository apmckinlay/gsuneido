## Language Introduction

Here is a simple example of Suneido code:

``` suneido
/* 
multi-line comment 
*/ 
calc = function (x, y = 0, dbl = false)
    {
    sum = x + y // single line comment
    if dbl is true
        sum *= 2
    return 'result is ' $ sum
    }
calc(123) => "result is 123"
calc(123, 456) => "result is 579"
calc(10, dbl:, y: 5) => "result is 30"
calc("123") => "result is 123"
```

Some features of the Suneido language:

-	C/C++/Java style comments, either `/*...*/` or `//...`
-	Functions are declared using the keyword `function`
-	Functions (and classes) are "values" - they can be assigned to variables etc.
-	Function arguments can have default values
-	Function calls can use named arguments, in which case order is not significant.  If a value isn't supplied for named argument, it defaults to true.
-	Semicolons are not required (but are allowed)
-	Parenthesis are not required around control (if, while, etc.) expressions (but are allowed)
-	Variables do not need to be declared, and may contain any type - i.e. Suneido is a dynamically typed language
-	[Local Variables](<Names/Local Variables.md>) start with a lower case letter. Global names start with an upper case letter.
-	Names starting with an underscore are 
	[Dynamic Variables](<Names/Dynamic Variables.md>) i.e. they are available to called functions without being explicitly passed.
-	You cannot assign to 
	[Global Names](<Names/Global Names.md>), they are either built-in or defined in libraries in the database. There is a global Suneido object which can be used to store global information.
-	`is`, `isnt`, `not`, `and`, `or` are normally used but `==`, `!=`, `!`, `&&`, `||` are allowed
-	Suneido includes all of the C/C++ operators, including assignment operators (e.g. *=), increment, decrement, and bit manipulation, and adds regular expression operators =~ and !~
-	Automatic conversions between numbers and strings.
-	String concatenation is done with the "$" operator, not by overloading "+".  This makes "123" + 456 and 123 $ "456" unambiguous.
-	Operators have consistent result types. e.g. "+" always produces a number, "$" always produces a string.
-	[String](<Basic Data Types/String.md>) literals can be written with either single or double quotes. This makes it easier to include quotes within strings.
-	Suneido includes the same control statements as C/C++ 
	[if-else](<Statements/if.md>), 
	[while](<Statements/while.md>), 
	[do-while](<Statements/do while.md>), 
	[switch](<Statements/switch.md>), but unlike C/C++ conditions must evaluate to true or false or you'll get an error.  
	[switch](<Statements/switch.md>) is also
	slightly different - you can't "fall through" cases, and cases allow multiple expressions rather than integer constants. Switches also throw an exception for unhandled values. There is also a version of the "for" statement that iterates through objects (see below).


Functions can take variable numbers of arguments:

``` suneido
max = function (@args)
    {
    max = args[0] // will throw exception if no arguments
    for x in args
        if x > max
            max = x
    return max
    }
max(3, 6, 2, 5) => 6 
max("joe", "fred", "sam", "mike") => "sam"
```

@args puts all the arguments into a Suneido "object" (see below).  for x in args iterates through all the values in an object, in this case all the arguments.

You can also pass a pre-assembled set of arguments, using @.

``` suneido
ages = #(23, 67, 34, 19)
max(@ages) => 67
```

Memory is garbage collected - there is no explicit allocation or freeing.

The basic data types in Suneido include: boolean (true and false), number, string, date, and object.

Suneido has a single numeric type - decimal floating point.  Keeping numbers in decimal rather than binary allows exact representation of decimals, e.g. for amounts of money.  Numbers have 16 digits of precision with an exponent range of plus or minus 512.

Strings are not null terminated so they can store arbitrary binary data as well as text.  Strings are immutable i.e. there is no way to "modify" an existing string.  (Objects are the only basic data type that is mutable.)  This means that substring (s[i,n]) does not need to do any copying - it just creates a new string that points to part of the old one.  For speed, concatenation is "lazy"; it creates a linked list of string segments, which are automatically combined when a single string is required.  This greatly reduces the amount of allocation and copying required to manipulate strings.

Suneido could be called a "pure" object-oriented language in the sense that all values (including literals) are "objects" that can have methods.  For example:

``` suneido
"hello world"[3::2] => "lo"
97.Chr() => "a"
```

However, unlike some object-oriented languages such as Smalltalk or Java, Suneido has standalone functions as well as methods in classes.

User defined methods can be added to the built-in classes by creating methods in specially named classes: Numbers, Strings, Dates, Objects, etc.

Suneido has only two "scopes", global and local.  Global names must be capitalized i.e. start with an upper case letter.  Globals are either built-in (to the executable) or user defined in
libraries in the database. Global names are not variables - they cannot be assigned to by code.  Names that start with a lower case letter are local to a single function.  Currently there are no
packages or modules with separate namespaces.

Like Java, Suneido compiles to a stack oriented "byte code" which is then interpreted.  For example:

``` suneido
function (x) { return x * 100 }.Disasm()
=>  push auto x
    push int 100
    *
    return
```

Linking is done dynamically at run time.  This means there are no restrictions on compiling code with calls to undefined functions or inheriting from undefined base classes.  Of course, if you try to actually access something undefined, you'll get a run time error.

**Objects**

Suneido has a single "universal" container type - objects.  Unlike other languages like Python or Ruby, you don't have to worry about what type of container to use.  Objects can be used as vectors (single dimensional arrays) or as *maps* (*dictionaries*) or both at the same time.  Internally, objects have a vector part and a hash map part. Classes and instances are similar but only have a map part.

``` suneido
x = Object() 
for (i = 0; i < 6; ++i) 
    x[i] = i // access as a vector 
x.Add("six") // same effect as x[6] = "six"
x => #(0, 1, 2, 3, 4, 5, "six")
x = Object(name: "fred", age: 25) 
x.married = true
x => #(name: "fred", age: 25, married: true) 
x.name => "fred"
m = "age" 
x[m] => 25 // same as x.age
```

**Exceptions**

Suneido has exception handling similar to C++ or Java, with try, catch, and throw.  However, currently, Suneido exceptions are strings rather than class instances.

``` suneido
try
    ...
catch (exception)
    ...
```

The catch portion is optional if you simply want to ignore exceptions:

``` suneido
try Database("drop mytable")
```

An uncaught exception calls a Handler function - defined in the standard library as the debugger.

You can throw your own exceptions:

``` suneido
if x < 0
    throw "square root: invalid negative argument: " $ x
```

**Blocks**

Suneido also includes Smalltalk style "blocks" (closures).  Basically, a block is a section of code within a function, that can be called like a function, but that operates within the context of the function call that created it (i.e. shares its local variables).

Blocks can be used to implement user defined "control constructs".  (In Smalltalk, all control constructs are implemented with blocks.)  For example, you could implement your own version of "for each":
<pre>
for_each = function (list, block)
    {
    for (i = 0; i &lt; list.Size(); ++i)
        block(list[i])
    }
list = #(12, 34, 56)
for_each(list)
    { |x| Print(x) }

<br />
 =&gt; 12
    34
    56
</pre>

Suneido treats a block immediately following a function call as an additional argument.

Blocks can also be used to execute sections of code in specific "contexts".  For example, the Catch function traps exceptions and returns them.  (This is useful in unit tests to verify that expected exceptions occur.)

``` suneido
catcher = function (block)
    {
    try
        return block()
    catch (x)
        return x
    }
catcher( { xyz } ) => "unitialized variable: xyz"
```

But the interesting part is that a block can outlive the function call that created it, and when it does so, it keeps its context (set of local variables). For example:

``` suneido
make_counter = function (next)
    { return { next++ } }
counter = make_counter(10)
Print(counter())
Print(counter())
Print(counter())
=>  10
    11
    12
```

In this example, make_counter returns a block. The block returns next++.  You see this type of code in Lisp / Scheme.

**Classes**

Classes are read-only objects containing methods (functions) and members (data).  For example:

``` suneido
counter = class
    {
    New(start = 0)
        { .count = start }
    Next()
        { return .count++ }
    }
ctr = new counter
Print(ctr.Next())
Print(ctr.Next())
Print(ctr.Next())
=>  0
    1
    2
```

Within a method, the special name "this" refers to the current object.  So members and methods of the current object can be referred to as "this.name" but this is normally shortened to simply ".name".

Member / method names are public if they start with an upper case letter and private if they start with a lower case letter.  Private names can only be accessed by methods of the class.  In the above example, "Next" is a public method and "count" is a private member.

New is the constructor.  (C++ and Java use the name of the class instead of "New", but this wouldn't work for Suneido because it allows "anonymous" classes i.e. with no name.)  Base class New's are automatically called first.  You can supply arguments to base class constructors by making the first statement of New a call to "super".  (Similar to Java.)

``` suneido
class
    {
    New(x, y, z)
        {
        super(x, y)
        ...
```

"super" can also be used to call base class versions of methods:

``` suneido
mystack = class : Stack
    {
    Push(x)
        {
        ...
        super.Push(x)
        ...
        }
```

Normally, "calling" a class is the same as creating an instance.  (This can be changed by defining a method called "CallClass".)  So the following are equivalent:

``` suneido
ctr = counter(5)
ctr = new counter(5)
```

Any method in Suneido can be a "class" or "static" method, providing it doesn't need to set any members.  If a method is called on a class rather than an instance, "this" will be the class itself, and therefore read- only.

If a nonexistent *member* is accessed, Suneido will call a method whose name is the member name prefixed by "Getter_" or "getter_", if present. This allows you to change a data member to a method without having to alter code that uses it. (Bertrand Meyer calls this "Uniform Access" - a client should be able to access a property of an object using a single notation, whether the property is implemented by memory or computation.)

If a nonexistent *method* is called, Suneido will call a method called "Default" if it is present.  This can be used to implement delegation, or to handle methods that are not pre-determined.

``` suneido
Stack.Members()
=> #("Pop", "Top", "Push")
```

Object syntax can also be used to access methods or members "indirectly":

``` suneido
m = "Add"
ob[m](value) // equivalent to ob.Add(value)
```

**"Missing" Features**

-	no declarations
-	no preprocessor i.e. #define or #include
-	no enums
-	no multiple inheritance
-	no operator overloading
-	no goto
-	no finally on try-catch
-	no explicit pointers
-	no explicit free'ing of memory
-	no separate module or package namespaces
-	no "protected" members / methods (just private and public)
-	only constant "static" class data members


**Summary**

The Suneido language is:

-	small, simple
-	object-oriented
-	dynamically typed - no declarations
-	functionally similar to Smalltalk
-	familiar to anyone who has programmed in C, C++, Java etc.
-	safe - automatic memory management (garbage collection) and no pointers