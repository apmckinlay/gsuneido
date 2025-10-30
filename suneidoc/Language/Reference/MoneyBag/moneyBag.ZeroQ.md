#### moneyBag.Zero?

``` suneido
() => true or false
```

Returns true if all the amounts in the money bag are zero, false otherwise.

For example:

``` suneido
mb = new MoneyBag
mb.Zero?()
    => true

mb = new MoneyBag
mb.Plus(100, 'CAD')
mb.Zero?()
    => false

mb = new MoneyBag
mb.Plus(100, 'CAD')
mb.Plus(-100, 'CAD')
mb.Zero?()
    => true
```