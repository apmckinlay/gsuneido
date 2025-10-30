### Rlog

``` suneido
Rlog(name, s, multi_line = false)
```

Writes to rotating log files, by default one per day, keeping the last 10.

The log file name is made from name $ number $ ".log". The number depends on the date.

Log entries are prefixed by the date and time.

Within the message %m is replaced with [MemoryArena](<MemoryArena.md>)() and %s by [Database.SessionId](<../../Database/Reference/Database/Database.SessionId.md>)()

For example:

``` suneido
Rlog("app", "%m start up")
```

would write to e.g. "app3.log"

``` suneido
20160103.120037556 47321088 startup
```
The multi_line option is intended for multi-line messages. For example:
``` suneido
Rlog("app", "hello\nworld", multi_line:)
```

would write:

``` suneido
20160103.120336868 -+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
hello
world
```

Logging is done via [ServerEval](<ServerEval.md>)() which means the file will normally be written to the same directory as suneido.db