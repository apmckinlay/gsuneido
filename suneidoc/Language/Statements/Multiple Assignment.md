### Multiple Assignment

``` suneido
var, var, ... = function-call
```

Assigns multiple return values to multiple local variables. **Note**: multiple return/assign is not supported on suneido.js

-	This is a statement, not an expression like single assignment. It does not have a value. (e.g. if it is the last statement in a block)
-	The number of variables must match the number of return values.


If some of the return values are not needed, they can be assigned to '_' (or 'unused'), for example:

``` suneido
a, _, c = Fn()
```

See also: [return](<return.md>)