### Logical

Logical operators always produce a boolean result (true or false).

<div class="table-half-width">

|             |             | 
|   :----     |    :----    |
| `and   &&`  | logical and | 
| <code>or   \|\|</code> | logical or  | 
| `not   !`   | logical not | 

</div>

"and",  "or", and "not" are more readable synonyms for "&&", "||", and  "!"
and reduce the risk of confusing "&" and "|" with "&&" and "||".

**Note:** "and" and "or" use "short circuit" evaluation. 
i.e. If the left hand value is sufficient to determine the result, 
then the right hand side will not be evaluated.
This allows expressions like:

``` suneido
if (ob.Member?("name") and ob.name is "Fred")
```

This will only evaluate "ob.name" if "name" is a member,
so you won't get an "unitialized member" error.