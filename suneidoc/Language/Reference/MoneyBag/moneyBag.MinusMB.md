#### moneyBag.MinusMB

``` suneido
(moneyBag) => this
```

Subtracts all currency amounts in moneyBag from the corresponding currency amounts in this.

for example:

``` suneido
mb = new MoneyBag
mb.Plus(100, 'USD')
mb.Plus(50, 'CAD')

mb2 = new MoneyBag
mb2.Plus(25, 'USD')
mb2.Plus(10, 'CAD')

mb.MinusMB(mb2)
    => MoneyBag(USD: 75 CAD: 40)
```