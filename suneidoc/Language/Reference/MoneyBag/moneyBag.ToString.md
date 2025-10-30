#### moneyBag.ToString

``` suneido
() => string
```

Returns a string representation of the MoneyBag object.

For example:

``` suneido
mb = new MoneyBag
mb.Plus(100, 'USD')
mb.Plus(100, 'CAD')
mb.ToString()
    => "MoneyBag(USD: 100 CAD: 100)"
```