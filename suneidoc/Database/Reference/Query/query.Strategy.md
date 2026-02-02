<div style="float:right"><span class="builtin">Builtin</span></div>

#### query.Strategy

``` suneido
(formatted = false) => string
```

Returns the strategy for the optimized query, showing index usage and costing information.

```
0.500x 500/1_000 10_000+20_000
   \     \   \      \     \
    \     \   \      \     variable cost
     \     \   \      fixed cost
      \     \   population
       \     nrows
        fraction
```

variable cost
: The part of the cost that depends on the number of rows you read. If you don't read any rows this would be zero.

fixed cost
: The part of the cost that is independent of how many rows you read. e.g. the cost of creating a temp index.

nrows
: The estimated number of rows that will be produced. This can be way off. **nrows** does **not** incorporate **frac**

population
: The number of rows that we are drawing from. For example, a Where on a key selects a single row (nrows = 1) from a population of the table size.

frac
: frac is the estimated fraction of the rows that will be read. Variable cost already incorporates frac. frac = 0 means only Lookup, else frac < 1 means Select

See also: [QueryStrategy](<../QueryStrategy.md>), [Query.Strategy1](<Query.Strategy1.md>)