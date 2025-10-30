<div style="float:right"><span class="builtin">Builtin</span></div>

### Database.Dump

``` suneido
(tablename = "", filename = "", publicKey = "")
```

Dumps the entire database or the specified table to a file. If client-server, this happens on the server.

If tablename is "" (or omitted) the entire database is dumped.

If filename is "" (or omitted) it defaults to database.su or tablename.su

If publicKey is supplied, the output of dump is encrypted, equivalent to [OpenPGP](<../../../Language/Reference/OpenPGP.md>).PublicEncrypt

Equivalent to the `-dump` [command line option](<../../../Introduction/Command Line Options.md>)

Dump'ed files are intended to be loaded with the `-load` [command line option](<../../../Introduction/Command Line Options.md>) or with [Database.Load](<Database.Load.md>)