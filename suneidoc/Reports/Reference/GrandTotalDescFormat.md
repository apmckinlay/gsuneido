### GrandTotalDescFormat

``` suneido
(item, skip = .16)
```

GrandTotalDescFormat is derived from [GrandTotalFormat](<GrandTotalFormat>). It displays and justifies the **item** to the right.

For example: 

``` suneido
(Query
    (_output 
        desc: (GrandTotalDesc (Text 'Grand Total') ), 
        amount: (GrandTotal grandtotal_amount) )
    )
```