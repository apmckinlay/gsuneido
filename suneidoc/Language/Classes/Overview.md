### Overview

In Suneido, classes are read-only collections of named members and methods. Classes can derive/inherit from other classes. Instances can be created from a class. Instance inherit the members and methods of their class. Instances cannot define methods - method lookup starts in the class. (You can store a function in a member, but you cannot call it like a method.) Instances can override data members from their class.

**Note:** In Suneido, classes are *values*. This means they can be assigned to variables, passed to other functions, etc.

Class literals are written as follows:

class:
<pre>
<b>class</b> { members }
<b>class</b> <i>Base</i> { members }
<b>class</b> : <i>Base</i> { members }
<i>Base</i> { members } // same as class : Base { members }
</pre>

member:

``` suneido
name: literal
name ( params ) { statements } // short for name: function (params)
```

Class members may be optionally separated by semicolons.

Instances of classes are created with:

``` suneido
new class
new class ( arguments )
```

**Note**: new should not be used with built-in classes such as Object or File.

Since this is also the default definition of CallClass, instances of most classes can also be created by "calling" the class:
<pre>
class ( <i>arguments</i> )
</pre>

For example, you could create a stack with:

``` suneido
new Stack()
```

or just:

``` suneido
Stack()
```

String concatenation ($) is allowed in string literal members. This concatenation is done at compile time. This primarily useful for formatting source code. For example:

``` suneido
class
    {
    message: "now is the time\n" $
             "for all good men\n"
    }
```

The following methods are available on classes:

-	[.Base](<../Reference/Object/object.Base.md>)
-	[.Base?](<../Reference/Object/object.Base?.md>)
-	[.Eval](<../Reference/Object/object.Eval.md>)
-	[.Eval2](<../Reference/Object/object.Eval2.md>)
-	[.Iter](<../Reference/Object/object.Iter.md>) - only iterates the members of the instance, not its classes
-	[.GetDefault](<../Reference/Object/object.GetDefault.md>) - includes inherited members
-	[.Member?](<../Reference/Object/object.Member?.md>) - includes inherited members
-	[.Members](<../Reference/Object/object.Members.md>) - includes inherited members with all:
-	[.Method?](<../Reference/Object/object.Method?.md>)
-	[.MethodClass](<../Reference/Object/object.MethodClass.md>)
-	[.Readonly?](<../Reference/Object/object.Readonly?.md>) - always true for classes
-	[.Size](<../Reference/Object/object.Size.md>)


Instances of classes also support:

-	[.Copy](<../Reference/Object/object.Copy.md>)
-	[.Delete](<../Reference/Object/object.Delete.md>)