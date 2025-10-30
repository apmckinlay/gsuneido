<div style="float:right"><span class="builtin">Builtin</span></div>

#### query.Tree

``` suneido
() => tree
```

QueryTree returns a syntax tree for an optimized query (i.e. the strategy). It may be rearranged. If you want the tree for the original query use [Query.Parse](<Query.Parse.md>)
See 
[Query.Parse](<Query.Parse.md>) for the basic properties prior to optimization..
The following properties will are available once the query has been optimized (e.g. from [transaction.Query](<../Transaction/transaction.Query.md>) or [WithQuery](<../WithQuery.md>)):
`nrows`
: The estimated number of rows that will be produced by this operation node.

`pop`
: The "population" that nrows will be drawn from.

`fast1`
: true if the operation node will produce at most a single row e.g. key() in a fast (i.e. indexed) way.

`frac`
: The estimated fraction of nrows that will be read.

`fixcost`
: The fixed cost for this operation node, regardless of how many rows are read. e.g. the cost of building a temp index.

`varcost`
: The variable cost for this operation node, incorporating frac.

The following properties will only be available if the query has also been executed:
`tget`
: The TSC count for this sub-tree.

`tgetself`
: The TSC count for just just this operation node, after subtracting the count from child nodes.