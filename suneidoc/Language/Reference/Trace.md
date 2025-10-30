<div style="float:right"><span class="builtin">Builtin</span></div>

### Trace

``` suneido
(string) or (number, block)
```

Trace(string) outputs the string.

A numeric argument turns tracing on or off.  Trace(0) turns tracing off.  One or more of the following values or'ed together can be used to enable tracing:
`TRACE.CONSOLE`
: Send trace output to a console window on Windows GUI (gsuneido.exe) or to stdout otherwise. (WARNING: Don't close the trace window - this will terminate Suneido.)

`TRACE.LOGFILE`
: Send trace output trace.log text file in the current directory.

`TRACE.CLIENTSERVER`
: Trace requests and responses between client and server. Only works if you are running client-server, not standalone.

`TRACE.DBMS`
: Traces certain dbms related activity.

`TRACE.QUERIES`
: Trace database queries.

`TRACE.QUERYOPT`
: Trace top level database query optimization.

`TRACE.JOINOPT`
: Trace join query optimization.

`TRACE.RECORDS - gSuneido only`
: Trace rules and observers. Can be helpful to debug rule issues.

`TRACE.SLOWQUERIES`
: Trace database queries that read more than 100 records per result record. (e.g. due to non-indexed where's)

If you don't specify either CONSOLE or LOGFILE (trace.log) then output will go to both.

The block form is recommended to ensure that tracing is turned off.