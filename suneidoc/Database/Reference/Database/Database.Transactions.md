<div style="float:right"><span class="builtin">Builtin</span></div>

#### Database.Transactions

``` suneido
() => object
```

Returns a list of the outstanding **update** transactions.

Note: This list should be empty when nothing is running.

For example:

``` suneido
Transaction(update:)
	{|unused|
	Database.Transactions()
	}
=> #(123)
```

If the database has been locked after detecting corruption, then Database.Transactions will return #(0)