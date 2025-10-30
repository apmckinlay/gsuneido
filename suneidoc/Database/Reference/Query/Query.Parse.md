<div style="float:right"><span class="builtin">Builtin</span></div>

#### Query.Parse

``` suneido
(query) => tree
```

Query.Parse parses a query string and returns a query tree. To parse code use [Suneido.Parse](<../../../Language/Reference/Suneido/Suneido.Parse.md>).

Query.Parse does not optimize the query. To get an optimized query tree (i.e. the strategy), use [query.Tree](<query.Tree.md>)

Query.Parse returns a query syntax tree. 

Tree nodes have "properties". For example:

``` suneido
q = Query.Parse("tables join columns")
q.type
=> "join"
```

All nodes have the following properties:
`type`
: The type of node: table, extend, rename, project, summarize, where, view, sort, join, leftjoin, times, union, minus, intersect. Optimizing may introduce other types of nodes: nothing, projectnone, tablelookup, tempindex

`nchild`
: The number of children this node has - 0, 1, or 2. If nchild is 1 the node will have a "source" property. If nchild is 2 the node will have "source1" and "source2" properties. These source properties will be sub-tree nodes.

`string`
: Returns a query string for just this operation node.

`String`
: Returns a query string for the sub-tree starting at this node.

If the query has been optimized, string and String will give the strategy.

#### table or view
`name`
: The table or view name.

#### where
`expr`
: The where expression. See 
[Suneido.Parse](<../../../Language/Reference/Suneido/Suneido.Parse.md>) for the details of expressions.