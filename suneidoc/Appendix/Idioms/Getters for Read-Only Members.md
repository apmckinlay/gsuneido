### Getters for Read-Only Members

``` suneido
Getter_Size()
    {
    return .size
    }
```

If you want to make a member public, but read-only, use a Getter_ function. Users of the member do not need to know it is using a Getter_ function - they can access it just as if it were a normal public member. For example:

``` suneido
thing.Size
```

**Warning:** This will not prevent assignment to the public member, which will then override the Getter_ method.

See also: [Getter Methods](<../../Language/Classes/Getter Methods.md>)