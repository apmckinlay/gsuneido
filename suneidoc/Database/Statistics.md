## Statistics

During queries, the most used table columns (in where's) are tracked.
This information is written to a file called busycols.txt
every two hours (or on exit if running less than two hours).

When the database is compacted, busycols.txt is read (and then removed).
During the compaction process, statistics are gathered for the top 100 table columns.
At the end, this information is written to the _stats_ table (in binary form).
These statistics are then used to help query optimization.
For informational purposes, the statistics can be viewed via the "dbstats" table.