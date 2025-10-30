#### moneyBag.Merge

``` suneido
(moneyBag) => this
```

Adds the amounts from another money bag to this one.

For example:

``` suneido
mb1 = new MoneyBag
mb1.Plus(100, 'CAD')
mb1.Plus(200, 'USD')
mb2 = new MoneyBag
mb2.Plus(300, 'USD')
mb2.Plus(400, 'YEN')
mb1.Merge(mb2)
return mb1.Amounts()
	=> #(USD: 500, YEN: 400, CAD: 100)
```