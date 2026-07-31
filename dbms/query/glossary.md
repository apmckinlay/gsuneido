# Glossary

ap
: usual name for a query operation "approach"

approach
: one of the possible ways to execute a query operation (aka "strategy")

by
: the usual name for the join columns of a Join/LeftJoin

cost
: the estimated cost to execute a query or subquery (type Cost)

disjoint
: a column whose values differ between the two sources of a
  Compatible operation (Union/Intersect/Minus),

fixcost
: the portion of the cost that is "fixed"
  (i.e. does not depend on the frac)

fixed
: columns with known values due to where or extend operations (type Fixed)

frac
: the fraction of a query operation output that we expect to read

index
: a physical access path on a table or subquery,
  unlike keys, indexes are physical, not logical

keys
: short for candidate keys

logical
: things we know "logically" or "theoretically"
  e.g. keys are logical

physical
: whether there is a physical access path (i.e. an index) for execution
  e.g. indexes are physical

q
: usual name for a Query variable

query
: usual name for a query string

req
: the usual name for a Require variable

sels
: the usual name for a Sels variable;
  column/value pairs passed to Lookup/Select

singleton
: a query operation that returns at most one row (logical)

varcost
: the variable cost that depends on the frac
