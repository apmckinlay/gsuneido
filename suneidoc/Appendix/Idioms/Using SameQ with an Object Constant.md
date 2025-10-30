### Using Same? with an Object Constant

Sometimes you need a special unique value that is distinguishable from all other values.

For example, Objects.Count counts occurrences of any value. But the value is optional so we need a default value that we can distinguish. We can't use the normal default values of false or 0 or "" because we might want to count those.

We can use a default value of something like #(0) and take advantage of the Suneido compiler's constant sharing.

``` suneido
Count(value = #(0))
	{
	if Same?(value, #(0))
		return .Size()
```

The compiler will use a single instance of #(0) for the Count function. This will be a different instance from #(0) anywhere else.

**Note**: We need to use Same? instead of "is" to distinguish this particular instance.

**Note**: we can't use #() because Suneido shares a single global read-only empty object.