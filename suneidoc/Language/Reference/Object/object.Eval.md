<div style="float:right"><span class="builtin">Builtin</span></div>

#### object.Eval

``` suneido
( fn, @args ) => value
```

Calls fn with the supplied arguments as if it were a method of this object. 
i.e. with *this* set to this object. Applicable to objects, records, classes, and instances.

Returns whatever the function returns.

fn can be anything callable: block, function, class, or instance.

Eval is required because a normal function call does not propogate *this*.
i.e. A called standalone function normally does **not** have access to instance variables.

For example:

``` suneido
context = Object(Name: "Fred");
fn = function () { return .Name; };
context.Eval(fn);
=> "Fred"
```

**Note:** Eval does not provide access to private members since these are class specific.

If the first argument to Eval is a [Method](<../../../Appendix/Glossary/Method.md>) the instance or class it is bound to will be ignored.

See also: [object.Eval2](<object.Eval2.md>)