### Dynamic Variables

Dynamic variables start with an underscore followed by a lower case letter. Dynamic variables set in the current function may be accessed by any functions it calls (directly or indirectly). Access to callee dynamic variables is read-only; if a new value is assigned to a dynamic variable it will not affect previous values.

**Note**: When a dynamic variable is set to a mutable value (e.g. an object) then the called code can modify the contents of the value. The original dynamic variable hasn't changed, it still points to the same object. This is similar to how passing arguments in [function calls](<../Expressions/Function Calls.md>) works.

**Note:** An underscore followed by a *capital* letter is a reference to the previous value of an overridden [global name](<Global Names.md>).

See also: 
[Implicit Dynamic Parameters](<../Functions/Implicit Dynamic Parameters.md>)