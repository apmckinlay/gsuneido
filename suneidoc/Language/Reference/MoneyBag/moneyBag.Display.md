#### moneyBag.Display

``` suneido
(format = '-###,###,###.##') => string
```

returns a string with formatted dollar amounts suitable for displaying to the user.

For example:

``` suneido
mb = new MoneyBag
mb.Plus(100, 'USD')
mb.Plus(125, 'CAD')
mb.Display()
    => "100.00 USD, 125.00 CAD"
```