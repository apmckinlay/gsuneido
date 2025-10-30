### MoneyBag


MoneyBag is used to easily accumulate totals in multiple currencies.

For example:

``` suneido
mb = new MoneyBag
mb.Plus(100, 'CDN')
mb.Minus(100, 'CDN')
mb.Plus(400, 'USD')
mb.Plus(56.78, 'USD')
mb.Amounts()
	=> #(CDN: 0, USD: 456.78)
```
|     |
| --- |
| [MoneyBag](<MoneyBag/MoneyBag.md>) |
| [moneyBag.Amounts](<MoneyBag/moneyBag.Amounts.md>) |
| [moneyBag.Display](<MoneyBag/moneyBag.Display.md>) |
| [moneyBag.Merge](<MoneyBag/moneyBag.Merge.md>) |
| [moneyBag.Minus](<MoneyBag/moneyBag.Minus.md>) |
| [moneyBag.MinusMB](<MoneyBag/moneyBag.MinusMB.md>) |
| [moneyBag.Plus](<MoneyBag/moneyBag.Plus.md>) |
| [moneyBag.ToString](<MoneyBag/moneyBag.ToString.md>) |
| [moneyBag.Zero?](<MoneyBag/moneyBag.Zero?.md>) |

