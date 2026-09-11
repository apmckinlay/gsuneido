<div style="float:right"><span class="builtin">Builtin</span></div>

#### TypeChecker.Annotate

``` suneido
(arguments :object, references = #(), config = #()) => object
```

Infers the types in each of the **arguments** sources like
[TypeChecker.Infer](<TypeChecker.Infer.md>), but instead of reporting the types separately, writes them back into the source as inline comments and param/return annotations.

`arguments`
: arguments are object in order of source starting with the base class and ending with the current class

For example:

``` suneido
src = 'class
	{
	Total(n)
		{
		sum = 0
		for (i = 0; i < n; ++i)
			sum += i
		return sum
		}
	}'
Print(TypeChecker.Annotate([[src: src, name: "Adder"]]).result[0])
=>
class
	{
	Total(n) :number
		{
		sum /* number */ = 0
		for (i /* number */ = 0; i /* number */ < n; ++i /* number */)
			sum /* number */ += i /* number */
		return sum /* number */
		}
	}
```

See also: [TypeChecker.Infer](<TypeChecker.Infer.md>)