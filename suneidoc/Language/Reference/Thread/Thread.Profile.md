<div style="float:right"><span class="builtin">Builtin</span></div>

#### Thread.Profile

``` suneido
(block) => data
```

Run the block, collecting profiling data on the current thread.

It is better to profile on a separate thread so it is not affected by other activity on the main UI thread.

**Note:** Only available on x86 / amd64. (as of 2024-08-12)

For example:

``` suneido
p = Thread.Profile({ SvcTest() })
p.Sort!(By(#self))[-10..].Each(Print);;
=>
#(calls: 136, total: 24464753, self: 24464753, name: "Transactions.QueryOutput /* stdlib block */")
#(calls: 57, total: 206258552, self: 24925734, name: "Svc.Put /* stdlib method */")
#(calls: 687, total: 91360870, name: "DoWithTran /* stdlib function */", self: 27232776)
#(calls: 18057, total: 34088203, self: 34088203, name: "ClassHelp.ClassHelp_nesting /* stdlib method */")
#(calls: 2, total: 76513751, self: 50262221, name: "SystemChanges.GetState /* stdlib method */")
#(calls: 26765, total: 53415929, self: 53415929, name: "ScannerWithContext.ScannerWithContext_skip? /* stdlib method */")
#(calls: 11, total: 444230057, self: 74300918, name: "ClassHelp.MethodRanges /* stdlib method */")
#(calls: 17997, total: 140146978, self: 86731049, name: "ScannerWithContext.ScannerWithContext_next /* stdlib method */")
#(calls: 228, total: 156537089, self: 94885544, name: "RetryTransaction.CallClass /* stdlib method */")
#(calls: 17997, total: 331998093, self: 191851115, name: "ScannerWithContext.Next /* stdlib method */")
```

The returned data is a list of objects, each containing information on a particular function or method.
`name`
: The name of the function or method.

`calls`
: The number of times it was called.

`self`
: Time spent in this function.

`total`
: Time spent in this function and the functions it calls.

self and total times are from the CPU time stamp counter. They are not in any particular units and should only be used as relative measurements. Other activity in the system may affect the values. It's a good idea to run the code multiple times to get a more accurate result.