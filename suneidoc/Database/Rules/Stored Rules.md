### Stored Rules

These are regular database table columns that also happen to have a rule.

For example, to assign a unique timestamp to new records:

``` suneido
Rule_timestamp
function ()
	{
	return Timestamp();
	}
```

Since this rule has no dependency, it will only be triggered once.

Another use of stored rules is where the value needs to also be maintained
by an update trigger.  
For example, currency exchange rates are looked up by a rule,
and a converted amount is calculated by another rule.
A summary is then maintained of the converted amount.
When the exchange rate table is modified,
an update trigger is required to update the exchange rate and converted amount
in order to keep the summary accurate.