## Crash Recovery

At start up, Suneido performs a quick check of the recent activity in the database. Since the database is append only, a crash usually affects the end of the database.

If a problem is found then Suneido will perform a full check of the database (equivalent to running -check from the command line) and automatically attempt to repair.

If -check or [Database.Check](<../Database/Reference/Database/Database.Check.md>) fail, this will also cause the next start up to do a full check.

If you run [Database.Check](<../Database/Reference/Database/Database.Check.md>) (or [Database.Dump](<../Database/Reference/Database/Database.Dump.md>) which does a check first) and corruption is detected then the database will be "locked", putting the system into a read-only mode. Transactions will appear to succeed, but nothing will be written to the database.

If corruption has been detected and the database is locked then the [Server Monitor](<../Tools/Server Monitor.md>) will show:

``` suneido
Database damage detected - operating in read-only mode
```