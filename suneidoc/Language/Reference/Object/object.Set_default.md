<div style="float:right"><span class="builtin">Builtin</span></div>

#### object.Set_default

``` suneido
() => this // remove default value
(value) => this // set default value
```

Set the default value returned for non-existent members. (Without a default, attempts to access non-existent members will cause an exception.)

Calling Set_default with no argument removes any default value.

For example:

``` suneido
ob = Object()
ob.name
    => ERROR: member not found: "name"

ob = Object().Set_default("")
ob.name
    => ""

ob = Object().Set_default() // remove default
ob.name
    => ERROR: member not found: "name"
```

When you want to sum values, it is often useful to set the default to 0 
so you don't have to check whether the type has been encountered.  
For example:

``` suneido
sums = Object().Set_default(0)
QueryApply("payments")
    { |x|
    sums[x.type] += x.amount
    }
```

**Note:** If the default value is an object, then it will be copied and assigned to the member. (Otherwise, an access that returns a default value will** not** create the member.) This is useful, for example, when you need to sum values based on more than one "dimension":

``` suneido
sums = Object().Set_default(Object().Set_default(0))
QueryApply("payments")
    { |x|
    sums[x.type1][x.type2] += x.amount
    }
```