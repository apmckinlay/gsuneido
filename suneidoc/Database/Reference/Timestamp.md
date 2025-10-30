<div style="float:right"><span class="builtin">Builtin</span></div>

### Timestamp

``` suneido
() => date
```

Returns a unique timestamp.  i.e. Timestamp will never return the same value twice.

This is useful for system generated table keys.

In client/server mode, Timestamp comes from the server to ensure unique values across multiple clients.

**Warning:** Timestamp uses the system time,
so if the system time is set back this may result in duplicate values.

**As Of April 2023**: using a new extended format so Timestamp can generate more than 1000 timestamps per second without the time getting ahead.  Since more are available, the server now returns batches of timestamps instead of just one at a time. This reduces the number of requests to the server. Batches are only valid for 1 second to keep the time in sync.

**NOTE:** Because of the batching, timestamps from different clients and the server may be slightly out of order when within the same second.