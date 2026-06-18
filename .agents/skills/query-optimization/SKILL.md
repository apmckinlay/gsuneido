---
name: query-optimization
description: Reference for Suneido query optimization architecture and design
---

# Query Optimization

This document describes the query optimization system in `dbms/query/`.

## Overview

Query optimization transforms a parsed query into an efficient execution plan through a multi-phase pipeline:

```
ParseQuery → Transform → Optimize → SetApproach → Execute
```

**Phase 1: Transform** (rule-based)
- Bottom-up tree rewriting
- Pushes operations down (e.g., Where before Join)
- Combines consecutive operations (e.g., two Wheres)
- Detects and eliminates impossible cases (e.g., Nothing sources)
- Not cost-based — applies when structurally possible

**Phase 2: Optimize** (cost-based)
- Recursive cost estimation
- Chooses best execution strategy for each operation
- Considers physical indexes, temp indexes, and alternative strategies
- Returns (fixcost, varcost, approach) tuple

**Phase 3: SetApproach**
- Locks in the chosen strategy
- Inserts TempIndex nodes where needed
- Recursively applies to child operations

## Query Operation Hierarchy

```
Query (interface)
├── Table          — physical database table (leaf)
├── Nothing        — empty result (leaf)
├── ProjectNone    — single empty row (leaf)
├── schemaTable    — system tables (leaf)
│   ├── Tables, Columns, Indexes, Views, History
├── Query1         — single-source operations
│   ├── Extend     — add computed columns
│   ├── Project    — select/project columns (also Remove)
│   ├── Rename     — rename columns
│   ├── Sort       — sort by columns
│   ├── Summarize  — group by + aggregate
│   ├── TempIndex  — temporary index (inserted by optimizer)
│   ├── Where      — filter rows
│   └── View       — pass-through for views
└── Query2         — two-source operations
    ├── Compatible — base for set operations
    │   ├── Union
    │   └── Compatible1
    │       ├── Intersect
    │       └── Minus
    └── joinLike   — base for joins
        ├── Times  — cross product
        └── joinBase
            ├── Join
            ├── LeftJoin
            └── SemiJoin
```

## How Each Operation Uses Index/Requirement

| Operation | Index Semantic | Recursion Pattern | Helpers Used |
|---|---|---|---|
| **Table** | ordered (prefix match) | leaf | — |
| **Nothing** | ignore | leaf | — |
| **ProjectNone** | ignore | source: nil | — |
| **schemaTable** | reject if non-nil | leaf | — |
| **Extend** | pass-through (filter ext cols) | source: filtered index | — |
| **Rename** | pass-through (reverse rename) | source: renamed index | — |
| **TempIndex** | pass-through | source: nil | — |
| **Query1** (default) | pass-through | source: same index | — |
| **Sort** | ordered (replaces with sort.order) | source: sort.order | — |
| **Project** | grouped (or pass-through if unique) | source: via bestGrouped | `bestGrouped` |
| **Summarize** | grouped / ordered / mixed | source: varies by strategy | `bestGrouped`, `min3` |
| **Where** | ordered (index selection) | source: baseline + per-index | `ordered`, `newBestIndex`, `WhereCost` |
| **Times** | pass-through to src1, nil to src2 | both sources | — |
| **Join/LeftJoin** | ordered (src1) + lookup/grouped (src2) | both sources | `bestOrdered`, `bestLookupIndex`, `bestGrouped` |
| **SemiJoin** | pass-through (src1) + grouped (src2) | src1 direct | `bestGrouped` |
| **Intersect** | pass-through (src1) + lookup (src2) | src1 direct | `bestLookupIndex` |
| **Minus** | pass-through (src1) + lookup (src2) | src1 direct | `bestLookupIndex` |
| **Union** | ordered (merge) or lookup | both sources | `bestLookupIndex`, `bestMergeIndexes`, `min3` |

### Operation Categories

**Leaf nodes** (no recursion):
- Table, Nothing, ProjectNone, schemaTable

**Pass-through** (delegate to source):
- Query1 (default), Rename, Extend, TempIndex

**Requirement creators** (create their own requirement):
- Sort: creates `ReqOrdered(sort.order)`
- Join/LeftJoin src2: creates `ReqLookup(by)` or `ReqGrouped(by)`
- SemiJoin src2: creates `ReqGrouped(by)`
- Intersect/Minus src2: creates `ReqLookup`

**Requirement consumers** (combine parent req with own req):
- Project: uses `MergeReq` to combine parent req with `ReqGrouped(p.columns)`
- Summarize: uses `MergeReq` to combine parent req with `ReqGrouped(su.by)`

**Two-source operations**:
- Join, LeftJoin, SemiJoin, Times, Union, Intersect, Minus

## The Require System (v2)

The v2 optimizer uses a richer requirement model than v1's flat `index []string`.

### Use Enum

```go
type Use int
const (
    ReqUnordered Use = iota  // no ordering needed
    ReqOrdered               // rows must come in column order
    ReqGrouped               // rows must be grouped by columns (any order)
    ReqLookup                // individual point lookups needed
    ReqConflict = -1         // incompatible requirements
)
```

### Require Struct

```go
type Require struct {
    use  Use
    cols []string
}
```

### When to Use Each Type

**ReqUnordered**:
- No specific access pattern required
- Used by: mapCost (Project, Summarize), Times src2

**ReqOrdered**:
- Rows must arrive in a specific column order
- Used by: Sort, Where (for range scans), Join src1 (when parent needs ordering)

**ReqGrouped**:
- Rows must be grouped by columns (order within group doesn't matter)
- Used by: Project (for dedup), Summarize (for grouping), Join src2 (to-many)

**ReqLookup**:
- Individual point lookups needed
- Used by: Join src2 (to-one), Intersect/Minus src2, Union (lookup strategy)

**ReqConflict**:
- Returned by `MergeReq` when requirements are incompatible
- Triggers fallback to alternative strategy (e.g., projMap instead of projSeq)

## MergeReq

`MergeReq` combines two requirements symmetrically. Used when an operation has its own requirement AND must satisfy the parent's requirement.

### Which Operations Need MergeReq

Only **Project** and **Summarize** need MergeReq:

| Operation | Incoming req | Own requirement | Merge? |
|---|---|---|---|
| **Project** (non-unique) | from parent | `ReqGrouped(p.columns)` for dedup | **YES** |
| **Summarize** (seqCost) | from parent | `ReqGrouped(su.by)` for grouping | **YES** |
| All others | — | — | No |

### Merge Rules

```go
func MergeReq(req1 Use, cols1 []string, req2 Use, cols2 []string) (Use, []string)
```

**Common cases**:
- `ReqUnordered + ReqGrouped(cols)` → `ReqGrouped(cols)`
- `ReqOrdered(parent_cols) + ReqGrouped(op_cols)` → `orderedPlusGrouped` → `ReqOrdered(merged)` or `ReqConflict`
- `ReqGrouped(parent_cols) + ReqGrouped(op_cols)` → `ReqGrouped` if equal, else `ReqConflict`
- `ReqOrdered + ReqLookup` → `ReqOrdered` if equal cols
- `ReqGrouped + ReqLookup` → `ReqLookup` if lookup starts with grouped
- `ReqLookup + ReqLookup` → `ReqLookup` if equal, else `ReqConflict`

**Conflict handling**:
- When `MergeReq` returns `ReqConflict`, fall back to alternative strategy
- Project: use `projMap` (hash-based dedup) instead of `projSeq`
- Summarize: use `sumMap` or `sumIdx` instead of `sumSeq`

## Helper Functions

### bestGrouped

```go
func bestGrouped(source Query, mode Mode, index []string, frac float64, cols []string) bestIndex
```

Finds the best index with `cols` (in any order) as a prefix, taking fixed into consideration.

**Used by**: Project, Summarize, Join (src2 when not to-one), SemiJoin (src2)

**Logic**:
- If `index` is nil: tries all source indexes plus using `cols` directly as temp index
- If `index` is non-nil: only tries that specific index
- Uses `grouped()` to check if index satisfies grouping requirement

### bestLookupIndex

```go
func bestLookupIndex(source Query, mode Mode, nrows int, frac float64, cols []string) bestIndex
```

Finds the best index for `nrows` point lookups.

**Used by**: Join (src2 when to-one), Intersect, Minus, Union (optLookup)

**Logic**:
- `cols` restricts candidates to those grouped by cols (Join to-one)
- `nil` cols allows any lookup-eligible index (Intersect, Minus, Union)
- Falls back to logical keys if no physical index qualifies
- Uses `LookupCost` to estimate cost

### bestOrdered

```go
func bestOrdered(q Query, order []string, mode Mode, frac float64, fixed Fixed) bestIndex
```

Finds the best index satisfying a required order.

**Used by**: Join/LeftJoin (via `optOrdered` for src1)

**Logic**:
- Iterates all source indexes
- Uses `ordered()` to check if index satisfies order requirement
- Returns cheapest index that satisfies order

### ordered / grouped

```go
func ordered(index []string, order []string, fixed Fixed) bool
func grouped(index []string, cols []string, nColsUnfixed int, fixed Fixed) bool
```

Check if an index satisfies an ordering or grouping requirement, taking fixed columns into account.

**ordered**: index must have `order` columns as a prefix (in order)
**grouped**: index must have `cols` columns as a prefix (in any order)

### LookupCost

```go
func LookupCost(q Query, mode Mode, index []string, nrows int, frac float64) (Cost, Cost)
```

Returns the cost of performing `nrows` lookups on a query using the specified index.

**Logic**:
- Calls `optimize` to get base cost
- If approach is `tempIndex`, uses fixed lookup cost (200 for single table, 400 otherwise)
- Otherwise uses `q.lookupCost()`
- Multiplies by `nrows`

## Optimization Patterns

### Pass-Through Operations

These operations delegate to their source, possibly transforming the requirement:

**Rename**: reverse-maps column names
```go
func (r *Rename) optimize2(mode Mode, req *Require, frac float64) (Cost, Cost, any) {
    cols := r.renameRev(req.cols)
    fixcost, varcost := Optimize2(r.source, mode, &Require{req.use, cols}, frac)
    return fixcost, varcost, nil
}
```

**Extend**: filters out extended columns if source is fastSingle
```go
func (e *Extend) optimize2(mode Mode, req *Require, frac float64) (Cost, Cost, any) {
    cols := req.cols
    if e.source.fastSingle() {
        cols = e.filterSourceIndex(cols)
    } else if !set.Disjoint(cols, e.cols) {
        return impossible, impossible, nil
    }
    fixcost, varcost := Optimize2(e.source, mode, &Require{req.use, cols}, frac)
    return fixcost, varcost, nil
}
```

### Requirement Creators

These operations create their own requirement for their source:

**Sort**: replaces incoming req with `ReqOrdered(sort.order)`
```go
func (sort *Sort) optimize2(mode Mode, req *Require, frac float64) (Cost, Cost, any) {
    assert.That(req.use == ReqUnordered)
    fixcost, varcost := Optimize2(sort.source, mode, &Require{ReqOrdered, sort.order}, frac)
    return fixcost, varcost, nil
}
```

**Join src2**: creates `ReqLookup(by)` (to-one) or `ReqGrouped(by)` (to-many)
```go
if jt.toOne() {
    best2 = bestLookupIndex(src2, mode, nrows1, lookupFrac, by)
} else {
    best2 = bestGrouped(src2, mode, nil, frac2, by)
}
```

### Requirement Consumers

These operations use `MergeReq` to combine parent req with their own:

**Project** (non-unique):
```go
mergedUse, mergedCols := MergeReq(req.use, req.cols, ReqGrouped, p.columns)
if mergedUse == ReqConflict {
    return p.mapCost2(mode, req, frac)  // fallback to hash-based dedup
}
mergedReq := &Require{mergedUse, mergedCols}
fixcost, varcost := Optimize2(p.source, mode, mergedReq, frac)
```

**Summarize** (seqCost):
```go
mergedUse, mergedCols := MergeReq(req.use, req.cols, ReqGrouped, su.by)
if mergedUse == ReqConflict {
    return impossible, impossible, nil  // let sumMap or sumIdx win
}
mergedReq := &Require{mergedUse, mergedCols}
fixcost, varcost := Optimize2(su.source, mode, mergedReq, frac)
```

### Two-Source Operations

These operations optimize both sources, often with different requirements:

**Join/LeftJoin**:
- src1: pass-through (uses parent req)
- src2: creates `ReqLookup(by)` or `ReqGrouped(by)`
- Tries both forward and reverse (Join only), picks cheapest

**Union**:
- With required index: must use merge (both sources need `ReqOrdered`)
- Without required index: tries merge vs lookup (forward and reverse)

**Intersect/Minus**:
- src1: pass-through (uses parent req)
- src2: creates `ReqLookup`

### Where: The Most Complex Operation

Where has multiple optimization paths:

1. **Baseline filter**: optimize source with parent req, apply filter
2. **Index selection** (if source is table):
   - Iterate all physical indexes
   - For each index satisfying parent order:
     - Compute `WhereCost` using index selectivity fractions
     - Consider index range, index filter, data filter
   - Pick cheapest index

**Singleton optimization**:
- If Where condition fixes a key, result is single row
- Use `emptyKey` indexes, `fastSingle` returns true
- Avoids temp indexes

## Cost Model

### fixcost vs varcost

- **fixcost**: fixed cost of setting up the operation (e.g., building temp index)
- **varcost**: variable cost proportional to rows read

### frac (fraction)

Estimated fraction of rows that will be read:
- `frac = 1`: reading all rows
- `frac < 1`: reading subset (e.g., first/last, select)
- `frac = 0`: only lookups

Used to weight fixcost vs varcost:
- High frac: varcost dominates
- Low frac: fixcost dominates

### TempIndex Decision

`optTempIndex` decides whether to build a temporary sorted index:

```go
// Try three strategies:
optTI(best, q, mode, nil, frac, nrows, factorNone)      // no index
optTI(best, q, mode, index, frac, nrows, factorAll)     // required index
optTI(best, q, mode, bestIndex, frac, nrows, factorPre) // best partial index
```

**Constants**:
- `factorNone = 256`: penalty for no index
- `factorAll = 105`: penalty for required index
- `factorPre = 110`: penalty for partial index

**Cost formula**:
```go
fixcost := srccost + ticostAdj + 1000
fixcost += 100 * len(index)  // prefer fewer fields
if nrows > 0 {
    fixcost += factor * nrows * log(nrows)
}
varcost := frac * nrows * 100
if !q.SingleTable() {
    varcost *= 2
}
```

## Caching

### Per-Node Caching

Each query operation node has its own cache (embedded via `cache` struct in `queryBase`). This is important because:

1. **Multiple options at each node**: During optimization, we often try multiple strategies at the same node (e.g., different indexes for a Table, forward vs reverse for Join, merge vs lookup for Union)

2. **Recursive optimization revisits nodes**: When a parent operation tries different strategies, it may call `Optimize` on the same child node multiple times with the same requirements

3. **Cache avoids recalculating subtrees**: The cache stores the result of optimizing a subtree for a given `(use, cols, frac)`, avoiding redundant work

**Example**: A Join operation tries both forward and reverse directions. Each direction optimizes both sources. If both directions use the same requirement for source1, the cache prevents recalculating source1's subtree.

### Cache Key

Cache is keyed by `(use, cols, frac)`:
```go
type cacheEntry struct {
    approach any
    use      Use
    index    []string
    fixcost  Cost
    varcost  Cost
    frac     float64
}
```

**Why include `use`**:
- Different uses can produce different approaches
- `ReqOrdered(a,b)` might use index `(a,b,c)`
- `ReqGrouped(a,b)` might use index `(b,a,d)` (any order allowed)
- `ReqLookup(a,b)` might have different costs than sequential access

**Why include `frac`**:
- Different fractions can lead to different strategies
- Low frac favors low fixcost (avoid temp indexes)
- High frac favors low varcost (use efficient access paths)

### Cache Lifecycle

1. **Populated during Optimize**: As optimization explores options, results are cached
2. **Cleared after SetApproach**: Once the final approach is chosen, cache is cleared
3. **Not shared between nodes**: Each node has independent cache

### When Caching Helps

- **Multiple strategies at same node**: Join trying forward/reverse, Union trying merge/lookup
- **Parent explores options**: Parent tries different requirements on child
- **Complex queries**: Deep query trees with many optimization paths

### When Caching Doesn't Help

- **Simple queries**: Single operation, no alternatives to explore
- **Different requirements**: Each `(use, cols, frac)` combination is unique
- **One-shot optimization**: No repeated calls to same node

## Testing

### TestOptimize

Existing test suite in `optimize_test.go`:
- 256 lines covering all operation types
- Uses real query patterns
- Expected strings show chosen strategy

### Testing v2

V2 may produce different (potentially better) strategies than v1:
- Richer requirement types leading to different index choices
- `MergeReq` producing different merged requirements
- v2 helpers finding different optimal indexes

**Test approach**: verify v2 produces **valid** strategies, not necessarily identical to v1

```go
func test2(t *testing.T, query string, expected string) {
    q := ParseQuery(query, testTran{}, nil)
    q = q.Transform()
    fixcost, varcost := Optimize2(q, ReadMode, reqUnordered, 1)
    q = SetApproach(q, nil, 1, testTran{})
    assert.T(t).Msg(query).This(String(q)).Like(expected)
}
```

### Common Test Patterns

- Leaf nodes: v1 and v2 should be identical
- Pass-through: v1 and v2 should be identical
- MergeReq operations: v2 may choose different strategies
- Two-source operations: v2 may choose different join/union strategies

## Common Issues and Debugging

### Tracing Optimization

Enable trace output:
```go
trace.Set(int(trace.QueryOpt))
```

Shows:
- Each optimization call with requirements
- Cost estimates
- Strategy choices

### Common Bugs

**Impossible costs**:
- Returned when no valid strategy exists
- Check: are requirements too restrictive?
- Check: are physical indexes missing?

**Wrong strategy chosen**:
- Check cost estimates (fixcost, varcost)
- Check `frac` value
- Check if caching is interfering

**TempIndex not inserted**:
- Check `tempIndexable(mode)` — false in CursorMode
- Check if physical index is cheaper
- Check `ticostAdj` value

### Performance Issues

**Too many temp indexes**:
- Increase `factorNone`, `factorAll`, `factorPre`
- Check if physical indexes are available
- Check if `frac` estimates are accurate

**Slow optimization**:
- Check cache hit rate
- Check if too many recursive calls
- Profile to find bottleneck operations

**Slow execution**:
- Check chosen strategy (use `Strategy(q)`)
- Check if better physical indexes would help
- Check if Transform phase is working correctly

## Key Files

- `query.go`: Query interface, Optimize, SetApproach, Setup
- `require.go`: Require struct, Use enum, MergeReq
- `cache.go`: cache implementation
- `best.go`: bestGrouped, bestLookupIndex, bestOrdered, ordered, grouped
- `optimize_test.go`: test suite
- Operation files: `table.go`, `project.go`, `summarize.go`, `where.go`, `join.go`, etc.
