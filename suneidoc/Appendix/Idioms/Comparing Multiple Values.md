### Comparing Multiple Values

For example, a comparison function to sort by two fields is tricky:

``` suneido
return x.age is y.age ? x.name < y.name : x.age < y.age
```

and it gets more complicated the more fields you want to sort by.

Suneido compares lists of values, so the easier way to write this is:

``` suneido
return [x.age, x.name] < [y.age, y.name]
```

This is easy to extend to additional fields.

You can also use [object.ProjectValues](<../../Language/Reference/Object/object.ProjectValues.md>)

``` suneido
return x.ProjectValues(fields) < y.ProjectValues(fields)
```

If fields was #(age, name) this would be equivalent to the other versions.

See also: [object.Sort!](<../../Language/Reference/Object/object.Sort!.md>),
[By](<../../Language/Reference/By.md>)