# Glossary

block
: in-line anonymous function written {|params| ... }. Compiled as a normal function if it is not a **closure**

closure
: a **block** that shares variables with its containing function

column
: string name of a database column (see also field)

field
: numeric index of a field, usually in a record (see also column)

hamt
: hash array mapped trie immutable persistent data structure used by info and meta

info
: faster changing part of metadata, stats and index information

meta
: database metadata, consisting of layered info and schema

persist
: save the database state to storage, normally done once per minute rather than after every commit to reduce write amplification

schema
: slower changing part of metadata

action
: a "query" that performs an action, insert, update, or delete (QueryDo)

admin
: a database administration command, e.g. create, alter, drop

singleton
: a query that selects at most one record e.g. where a key equals a value
