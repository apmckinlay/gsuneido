# Glossary

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
