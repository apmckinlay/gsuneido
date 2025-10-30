## Using LruCache

**Category:** Coding

**Problem**

Your code is doing duplicate calculations or database queries that are slowing it down.

**Ingredients**

LruCache

**Recipe**

For example, if we wanted to total up transactions with a discount based on the customer:

``` suneido
total = 0
QueryApply('transactions')
    { |x|
    discount = Query1('customers where id = ' $ Display(x.id)).discount
    total += x.amount - x.amount * discount
    }
```

If there are a lot of transactions for only a small number of customers, this code is going to spend a lot of time looking up the same customer dicount over and over. We could speed it up a lot by saving the customer discounts in memory. But potentially, there could be too many customers to want to keep them all in memory.

LruCache implements a fixed size "cache". When the cache is full and it needs to make room for a new value, it throws out the least recently used value (LRU).

The first step is to extract the function (or block) that you want to cache:

``` suneido
get_discount = function (id)
    { return Query1('customers where id = ' $ Display(id)).discount }
total = 0
QueryApply('transactions')
    { |x|
    total += x.amount - x.amount * get_discount(x.id)
    }
```

Now we can easily introduce LruCache:

``` suneido
get_discount = function (id)
    { return Query1('customers where id = ' $ Display(id)).discount }
discounts = LruCache(get_discount, 100)
total = 0
QueryApply('transactions')
    { |x|
    total += x.amount - x.amount * discounts.Get(x.id)
    }
```

That's all there is to it. If the transactions involve less than 100 customers, all their discounts will be kept in memory. If there are more than 100 customers, the LruCache will keep the last 100.

If you want to see the difference, run Trace(TRACE.QUERY) and then run the code with and without LruCache.

If you want to test this code you'll need the tables. You can create them from QueryView with:

``` suneido
create customers (id, discount) key(id)

create transactions (t, id, amount) key(t) index(id) in customers
```

You can then enter some test data with IDE > Access a Query or Browse a Query.

**Tip**

If you want to find examples of where LruCache is used, a good way is to go to LruCache in LibraryView. You can do this using Find in Folders in LibraryView, or you can type LruCache on the WorkSpace and the hit F12 (or right click and choose Go To Definition). Once you're on the LruCache record in LibraryView, click on the L on the toolbar (or choose Tools > Show Locations). This will show you where LruCache is defined (useful if it's defined in several libraries) and where it's used. You can then choose one to go to it in LibraryView.

**Warning**

Premature optimization is a common mistake. Write the code as simply as possible. Time it on realistic data (not just a few test items). If it's too slow, then think about optimizing. And after making what you think is an improvement - time it again. You may find your wonderful idea doesn't help, in which case, go back to the simple version.