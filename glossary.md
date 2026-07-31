# Glossary

block
: in-line anonymous function written {|params| ... }. Compiled as a normal function if it is not a **closure**

closure
: a **block** that shares variables with its containing function

database
: append-only immutable database with MVCC

dbms
: layer above the database that provides the query language and client-server

hamt
: hash array mapped trie immutable persistent data structure used by info and meta

info
: faster changing part of metadata, stats and index information

meta
: database metadata, consisting of layered info and schema

persist
: save the database state to storage, normally done once per minute rather than after every commit to reduce write amplification

sb
: usual name for a strings.Builder variable

schema
: slower changing part of metadata
