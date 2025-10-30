### TotalDescFormat

``` suneido
(item, skip = .16)
```

TotalDescFormat is derived from [TotalFormat](<TotalFormat>). It displays and justifies the **item** to the right. 

For example: 

``` suneido
(Query
    (_output 
        desc: (TotalDesc (Text 'Total') ), 
        amount: (Total total_amount) )
    )
```