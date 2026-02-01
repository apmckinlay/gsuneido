---
name: suneido-database
description: Reference for Suneido database query language
---

# Suneido Database Query Language

## Query Syntax

*request* =
    query
    query **sort** [ **reverse** ] column [ , ... ]

*query* =
    table
    history(table)
    query **where** expression
    query **project** columns
    query **remove** columns
    query **join** query
    query **times** query
    query **union** query
    query **intersect** query
    query **minus** query
    query **rename** column **to** column [ , ... ]
    query **extend** column [ **=** expression ] [ , ... ]
    query **summarize** columns, [ column **=** ] function column [ , ... ]
    ( query )

## Update Syntax

*update* =
    **insert** { *column*: *value* [ , ... ] } **into** *query*
    **insert** *query* **into** *table*
    **update** *query* **set** *column* **=** *expression* [ , ... ]
    **delete** *query*

## Summarize Functions

Usage: `query summarize [by-columns,] [column =] function column`

Functions:
- `max`: Maximum value
- `min`: Minimum value
- `total`: Sum of values
- `average`: Average of values
- `count`: Count of values
- `list`: List of values

Note: If no name is specified, the result column is named `function_column` (e.g. `total_quantity`).

## Joins

### join (Natural Join)
`query join [ by (fields) ] query`
- Combines rows with equal common columns.
- Removes duplicates.
- Input queries must have at least one column in common.
- Does NOT output rows from the first table without a match.

### leftjoin (Left Outer Join)
`query leftjoin [ by (fields) ] query`
- Like join, but includes rows from the first table that have no match (with empty values for second table columns).

Note: `by (fields)` is an assertion only. Use `rename` to control join fields.
