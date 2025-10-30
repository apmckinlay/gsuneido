#### object.SumWith

``` suneido
(block) => number
```

Totals the results of the block applied to each value in the object.

For example:

``` suneido
data = #(
    (name: 'Fred', count: 456)
    (name: 'Joe', count: 123))
data.SumWith{ it.count }
    => 579
```

See also: [object.Sum](<object.Sum.md>)