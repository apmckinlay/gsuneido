### Stored Dependencies

Dependencies for fields are stored in fields with "_deps" suffixes. 
Currently, **these fields must be explicitly created.**

**Note:** These fields are only required if:

-	the field is a stored rule
-	the rule references other fields (i.e. has dependencies)


For example:

``` suneido
create transaction (date item price price_deps) key(date)

Rule_price
function ()
    {
    return GetItemPrice(.item);
    }
```